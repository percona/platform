package starlark

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
)

type Type string

type Check struct {
	Version uint32
	Type    Type
	Query   string
	Script  string
}

func Run(name, source, funcName string, input []map[string]interface{}) (*Check, error) {
	if !isInputValid(input) {
		return nil, errors.New("No valid input data")
	}

	resolve.AllowFloat = true
	resolve.AllowSet = true

	thread := &starlark.Thread{
		Name: name,
	}

	globals, err := starlark.ExecFile(thread, source, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "ExecFile")
	}

	if globals[funcName] == nil {
		return nil, errors.New("Function doesnt exists")
	}

	v, err := starlark.Call(thread, globals[funcName], starlark.Tuple{prepareRows(&input)}, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Call")
	}

	fmt.Println("Return type:" + reflect.TypeOf(v).String())

	return prepareResult(v)
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

func prepareResult(v starlark.Value) (*Check, error) {
	sd, ok := v.(*starlark.Dict)
	if !ok {
		return nil, errors.Errorf("expected *starlark.Dict, got %T", v)
	}

	res := make(map[string]interface{}, sd.Len())
	for _, tu := range sd.Items() {
		k := tu[0].(starlark.String).GoString()
		res[k] = starlarkToGo(tu[1])
	}

	return nil, nil
}
