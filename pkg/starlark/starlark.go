// Package starlark is executor for starklark.
package starlark

import (
	"github.com/pkg/errors"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"

	"github.com/percona-platform/platform/pkg/check"
)

// Run for execute starlark script.
func Run(name, script, funcName string, input []map[string]interface{}) (res *check.Result, err error) {
	res = new(check.Result)
	if !isInputValid(input) {
		res.Status = check.Fail
		res.Message = "No valid input data"
		return res, nil
	}

	resolve.AllowFloat = true
	resolve.AllowSet = true

	thread := &starlark.Thread{
		Name: name,
	}

	globals, err := starlark.ExecFile(thread, script, nil, nil)
	if err != nil {
		res.Status = check.Fail
		res.Message = "ExecFile: " + err.Error()
		return res, nil
	}

	if globals[funcName] == nil {
		res.Status = check.Fail
		res.Message = "Function doesnt exists"
		return res, nil
	}

	rows, err := prepareRows(&input)
	if err != nil {
		res.Status = check.Fail
		res.Message = err.Error()
		return res, nil
	}

	v, err := starlark.Call(thread, globals[funcName], starlark.Tuple{rows}, nil)
	if err != nil {
		res.Status = check.Fail
		res.Message = "Call: " + err.Error()
		return res, nil
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

func isInputValid(input []map[string]interface{}) bool {
	for _, item := range input {
		for _, v := range item {
			switch v.(type) {
			case uint64, int64, float64, string:
				continue
			default:
				return false
			}
		}
	}

	return true
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
