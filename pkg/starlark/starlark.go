// Package starlark is executor for starklark.
package starlark

import (
	"github.com/pkg/errors"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"

	"github.com/percona-platform/platform/pkg/check"
)

// Execute for execute starlark script.
func Execute(script string, input []map[string]interface{}) {
	name := "main"
	funcName := "main"

	run(name, script, funcName, input)
}

func run(name, script, funcName string, input []map[string]interface{}) (*check.Result, error) {
	res := new(check.Result)
	resolve.AllowFloat = true
	resolve.AllowSet = true

	thread := &starlark.Thread{
		Name: name,
	}

	globals, err := starlark.ExecFile(thread, script, nil, nil)
	if err != nil {
		res.Status = check.Fail
		res.Message = "ExecFile: " + err.Error()
		return res, errors.Wrap(err, "ExecFile: ")
	}

	if globals[funcName] == nil {
		res.Status = check.Fail
		res.Message = "Function doesnt exists"
		return res, errors.Errorf("Function %s doesnt exists", funcName)
	}

	rows, err := prepareRows(&input)
	if err != nil {
		res.Status = check.Fail
		res.Message = err.Error()
		return res, err
	}

	v, err := starlark.Call(thread, globals[funcName], starlark.Tuple{rows}, nil)
	if err != nil {
		res.Status = check.Fail
		res.Message = "Call: " + err.Error()
		return res, errors.Wrap(err, "Call: ")
	}

	switch v := v.(type) {
	case *starlark.Dict:
		for _, tu := range v.Items() {
			k := tu[0].(starlark.String).GoString()
			if k == "error" {
				res.Status = check.Fail
				res.Message = "Starlark script failed"

				return res, errors.New(string(tu[1].(starlark.String)))
			}
		}
		res.Status = check.Success
	default:
		res.Status = check.Success
	}

	return res, nil
}

func prepareRows(input *[]map[string]interface{}) (starlark.Tuple, error) {
	rows := make(starlark.Tuple, len(*input))
	for i, v := range *input {
		sv, err := goToStarlark(v)
		if err != nil {
			return nil, err
		}
		rows[i] = sv
	}
	rows.Freeze()

	return rows, nil
}
