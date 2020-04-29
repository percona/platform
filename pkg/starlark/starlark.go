// Package starlark is executor for starklark.
package starlark

import (
	"fmt"

	"github.com/pkg/errors"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"

	"github.com/percona-platform/platform/pkg/check"
)

// Env represents Starlark execution environment.
type Env struct {
	// Print is the client-supplied implementation of the Starlark 'print' function.
	// https://github.com/google/starlark-go/blob/master/doc/spec.md#print
	Print func(msg string)

	p *starlark.Program
}

// NewEnv creates a new Starlark execution environment.
func NewEnv(name, script string) (*Env, error) {
	predeclared := starlark.StringDict{}
	predeclared.Freeze()

	_, p, err := starlark.SourceProgram(name, script, predeclared.Has)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse script")
	}

	return &Env{
		p: p,
	}, nil
}

// noopPrint is a no-op 'print' implementation.
// It is a global function for a minor optimization (inlining, avoiding a closure).
func noopPrint(*starlark.Thread, string) {}

// print is a 'print' implementation that calls env.Print.
// It is a method for a minor optimization (avoiding a closure).
func (env *Env) print(t *starlark.Thread, msg string) {
	f := t.CallFrame(1)
	env.Print(fmt.Sprintf("%s %s %s: %s\n", t.Name, f.Pos.String(), f.Name, msg))
}

// run executes function with a given name with given arguments and returns result and fatal error.
// threadName is used only for debugging.
func (env *Env) run(threadName, funcName string, args starlark.Tuple) (starlark.Value, error) {
	thread := &starlark.Thread{
		Name:  threadName,
		Print: noopPrint,
	}
	if env.Print != nil {
		thread.Print = env.print
	}

	globals, err := env.p.Init(thread, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: failed to init script", threadName)
	}
	globals.Freeze()

	fn := globals[funcName]
	if fn == nil {
		return nil, errors.Errorf("%s: function %s is not defined", threadName, funcName)
	}

	v, err := starlark.Call(thread, fn, args, nil)
	if err != nil {
		if ee, ok := err.(*starlark.EvalError); ok {
			return nil, errors.Wrapf(err, "%s: failed to execute function %s\n%s", threadName, funcName, ee.CallStack.String())
		}
		return nil, errors.Wrapf(err, "%s: failed to execute function %s", threadName, funcName)
	}

	v.Freeze()
	return v, nil
}

// Run executes function 'check' with given query results.
// Id is used to separate that execution from other and used only for debugging.
func (env *Env) Run(id string, input []map[string]interface{}) (*check.Result, error) {
	rows, err := prepareInput(input)
	if err != nil {
		return nil, err
	}

	v, err := env.run(id, "check", starlark.Tuple{rows})
	if err != nil {
		return nil, err
	}

	switch v := v.(type) {
	case *starlark.Dict:
		// TODO https://jira.percona.com/browse/SAAS-84
		return &check.Result{
			Status:  "status",
			Message: "message",
		}, nil
	default:
		return nil, errors.Errorf("unhandled result type %T", v)
	}
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

// modify unavoidable global state once on package initialization to avoid race conditions
//nolint:gochecknoinits
func init() {
	resolve.AllowFloat = true
	resolve.AllowSet = true
}
