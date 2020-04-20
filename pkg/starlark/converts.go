package starlark

import (
	"reflect"

	"github.com/pkg/errors"
	"go.starlark.net/starlark"
)

func goToStarlark(v interface{}) (starlark.Value, error) {
	switch v := v.(type) {
	case int:
		return starlark.MakeInt(v), nil
	case string:
		return starlark.String(v), nil
	case float32:
		return starlark.Float(float64(v)), nil
	case float64:
		return starlark.Float(v), nil
	case bool:
		return starlark.Bool(v), nil
	default:
		return nil, errors.New("goToStarlark: Unhandled type " + reflect.TypeOf(v).String())
	}
}
