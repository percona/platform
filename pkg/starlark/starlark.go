// Package starlark provides Starlark execution environment.
package starlark

import (
	"fmt"

	"github.com/pkg/errors"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"

	"github.com/percona-platform/platform/pkg/check"
)

// Recover from panics in production code (we don't want all PMMs to crash after SaaS update),
// but crash in tests and during fuzzing.
// TODO Remove completely once Starlark is running in a separete process: https://jira.percona.com/browse/SAAS-63
//nolint:gochecknoglobals
var doRecover = true

// PrintFunc represents fmt.Println-like function that is used by Starlark 'print' function implementation.
type PrintFunc func(args ...interface{})

// Env represents Starlark execution environment.
type Env struct {
	p *starlark.Program
}

// NewEnv creates a new Starlark execution environment.
func NewEnv(name, script string) (env *Env, err error) {
	if doRecover {
		defer func() {
			if r := recover(); r != nil {
				err = errors.Errorf("%v", r)
			}
		}()
	}

	predeclared := starlark.StringDict{}
	predeclared.Freeze()

	var p *starlark.Program
	_, p, err = starlark.SourceProgram(name, script, predeclared.Has)
	if err != nil {
		err = errors.Wrap(err, "failed to parse script")
		return
	}

	env = &Env{
		p: p,
	}
	return
}

// noopPrint is a no-op 'print' implementation.
// It is a global function for a minor optimization (inlining, avoiding a closure).
func noopPrint(*starlark.Thread, string) {}

// run executes function with a given name with given arguments and returns result and fatal error.
// threadName is used only for debugging.
// print is a user-suplied function for Starlark 'print'.
func (env *Env) run(funcName string, args starlark.Tuple, threadName string, print PrintFunc) (starlark.Value, error) {
	thread := &starlark.Thread{
		Name:  threadName,
		Print: noopPrint,
	}
	if print != nil {
		thread.Print = func(t *starlark.Thread, msg string) {
			// make it look similar to starlark.CallStack.String
			fr := t.CallFrame(1)
			print("["+t.Name+"]", fr.Pos.String()+":", "in", fr.Name+":", msg)
		}
	}

	predeclared := starlark.StringDict{}
	predeclared.Freeze()

	globals, err := env.p.Init(thread, predeclared)
	if err != nil {
		if ee, ok := err.(*starlark.EvalError); ok {
			// tweak message, but keep original type, callstack, and cause
			ee.Msg = fmt.Sprintf("[%s] failed to init script: %s\n%s", threadName, ee.Msg, ee.CallStack)
			return nil, ee
		}
		return nil, errors.Wrapf(err, "[%s] failed to init script", threadName)
	}
	globals.Freeze()

	fn := globals[funcName]
	if fn == nil {
		return nil, errors.Errorf("[%s] function %s is not defined", threadName, funcName)
	}

	v, err := starlark.Call(thread, fn, args, nil)
	if err != nil {
		if ee, ok := err.(*starlark.EvalError); ok {
			// tweak message, but keep original type, callstack, and cause
			ee.Msg = fmt.Sprintf("[%s] failed to execute function %s: %s\n%s", threadName, funcName, ee.Msg, ee.CallStack)
			return nil, ee
		}
		return nil, errors.Wrapf(err, "[%s]: failed to execute function %s", threadName, funcName)
	}

	v.Freeze()
	return v, nil
}

// Run executes function 'check' with given query results.
// Id is used to separate that execution from other and used only for debugging.
// print is a user-suplied Starlark 'print' function implementation.
func (env *Env) Run(id string, input []map[string]interface{}, print PrintFunc) (res []check.Result, err error) {
	if doRecover {
		defer func() {
			if r := recover(); r != nil {
				err = errors.Errorf("%v", r)
			}
		}()
	}

	var rows *starlark.List
	rows, err = prepareInput(input)
	if err != nil {
		return
	}

	var output starlark.Value
	output, err = env.run("check", starlark.Tuple{rows}, id, print)
	if err != nil {
		return
	}

	res, err = parseScriptOutput(output)
	return
}

func prepareInput(input []map[string]interface{}) (*starlark.List, error) {
	values := make([]starlark.Value, len(input))
	for i, v := range input {
		sv, err := goToStarlark(v)
		if err != nil {
			return nil, err
		}
		values[i] = sv
	}

	l := starlark.NewList(values)
	l.Freeze()
	return l, nil
}

// parseScriptOutput parses and validates script output and returns slice of Results.
func parseScriptOutput(v starlark.Value) ([]check.Result, error) {
	switch v := v.(type) {
	case starlark.Tuple:
		if v.Len() != 2 {
			return nil, errors.New("script has invalid output")
		}

		errMsg, err := parseErrorMessage(v.Index(1))
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse error message")
		}

		if errMsg != "" {
			return nil, errors.Errorf("script error: %s", errMsg)
		}

		results, err := parseResults(v.Index(0))
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse results list")
		}

		for _, result := range results {
			if err = result.Validate(); err != nil {
				return nil, err
			}
		}

		return results, nil
	default:
		return nil, errors.Errorf("unhandled result type %T", v)
	}
}

// parseResults returns slice of results parsed from starlark value.
func parseResults(v starlark.Value) ([]check.Result, error) {
	val, err := starlarkToGo(v)
	if err != nil {
		return nil, err
	}

	rs, ok := val.([]interface{})
	if !ok {
		return nil, errors.Errorf("results list has wrong type: %T", val)
	}

	results := make([]check.Result, len(rs))
	for i, r := range rs {
		m := r.(map[string]interface{})
		var sum, desc string
		var sev check.Severity

		if v, ok := m["summary"]; ok {
			if sum, ok = v.(string); !ok {
				return nil, errors.Errorf("summary field has wrong type: %T", v)
			}
		}

		if v, ok := m["description"]; ok {
			if desc, ok = v.(string); !ok {
				return nil, errors.Errorf("description field has wrong type: %T", v)
			}
		}

		if v, ok := m["severity"]; ok {
			sevS, ok := v.(string)
			if !ok {
				return nil, errors.Errorf("severity field has wrong type: %T", v)
			}

			sev = check.StrToSeverity(sevS)
		}

		results[i] = check.Result{
			Summary:     sum,
			Description: desc,
			Severity:    sev,
		}
	}

	return results, nil
}

// parseErrorMessage returns error message parsed from starlark value.
func parseErrorMessage(v starlark.Value) (string, error) {
	val, err := starlarkToGo(v)
	if err != nil {
		return "", err
	}

	str, ok := val.(string)
	if !ok {
		return "", errors.Errorf("error message has wrong type: %T", val)
	}

	return str, nil
}

// modify unavoidable global state once on package initialization to avoid race conditions
//nolint:gochecknoinits
func init() {
	resolve.AllowFloat = true
	resolve.AllowSet = true
}
