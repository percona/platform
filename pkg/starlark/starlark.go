// Package starlark is executor for starklark
package starlark

import (
	"github.com/percona-platform/platform/pkg/check"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
)

// Function for execute starlark script
func Run(name, script, funcName string, input []map[string]interface{}) (res *check.Result) {
	res = new(check.Result)
	if !isInputValid(input) {
		res.Status = check.Fail
		res.Message = "No valid input data"
		return res
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
		return res
	}

	if globals[funcName] == nil {
		res.Status = check.Fail
		res.Message = "Function doesnt exists"
		return res
	}

	_, err = starlark.Call(thread, globals[funcName], starlark.Tuple{prepareRows(&input)}, nil)
	if err != nil {
		res.Status = check.Fail
		res.Message = "Call: " + err.Error()
		return res
	}

	res.Status = check.Success
	return res
}

func isInputValid(input []map[string]interface{}) bool {
	for _, item := range input {
		for _, v := range item {
			switch v.(type) {
			case int, float32, float64, string:
				continue
			default:
				return false
			}
		}
	}

	return true
}

func prepareRows(input *[]map[string]interface{}) starlark.Tuple {
	rows := make(starlark.Tuple, len(*input))
	for i, m := range *input {
		sd := starlark.NewDict(len(m))
		for k, v := range m {
			sv := goToStarlark(v)
			if err := sd.SetKey(starlark.String(k), sv); err != nil {
				return rows
			}
		}
		rows[i] = sd
	}
	rows.Freeze()

	return rows
}
