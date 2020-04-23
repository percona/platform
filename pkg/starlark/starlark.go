// Package starlark is executor for starklark.
package starlark

import (
	"github.com/pkg/errors"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"

	"github.com/percona-platform/platform/pkg/check"
)

// Execute for execute starlark script.
func Execute(script string, input []map[string]interface{}) (*check.Result, error) {
	return run("check", script, "check", input)
}

func run(name, script, funcName string, input []map[string]interface{}) (*check.Result, error) {
	thread := &starlark.Thread{
		Name: name,
	}

	globals, err := starlark.ExecFile(thread, script, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "ExecFile: ")
	}

	if globals[funcName] == nil {
		return nil, errors.Errorf("Function %s doesnt exists", funcName)
	}

	rows, err := prepareRows(input)
	if err != nil {
		return nil, err
	}

	v, err := starlark.Call(thread, globals[funcName], starlark.Tuple{rows}, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Call: ")
	}

	switch v := v.(type) {
	case *starlark.Dict:
		for _, tu := range v.Items() {
			k := tu[0].(starlark.String).GoString()
			if k == "error" {
				return nil, errors.New(string(tu[1].(starlark.String)))
			}
		}
		return nil, nil
	default:
		return nil, nil
	}
}

func prepareRows(input []map[string]interface{}) (starlark.Tuple, error) {
	rows := make(starlark.Tuple, len(input))
	for i, v := range input {
		sv, err := goToStarlark(v)
		if err != nil {
			return nil, err
		}
		rows[i] = sv
	}
	rows.Freeze()

	return rows, nil
}

// modify unavoidable global state once on package initialization to avoid race conditions
//nolint:gochecknoinits
func init() {
	resolve.AllowFloat = true
	resolve.AllowSet = true
}
