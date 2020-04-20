package starlark

import (
	"reflect"

	"github.com/pkg/errors"
	"go.starlark.net/starlark"
)

func goToStarlark(v interface{}) (starlark.Value, error) {
	switch v := v.(type) {
	case uint64:
		return starlark.MakeInt(int(v)), nil
	case int64:
		return starlark.MakeInt(int(v)), nil
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

func starlarkToGo(v starlark.Value) (interface{}, error) {
	switch v := v.(type) {
	case starlark.Bool:
		return bool(v), nil
	case starlark.Int:
		if i, ok := v.Int64(); ok {
			return int64(i), nil
		}
		if u, ok := v.Uint64(); ok {
			return uint64(u), nil
		}
		return nil, errors.New("starlarkToGo: Unhandled type " + reflect.TypeOf(v).String())
	case starlark.Float:
		return float64(v), nil
	case starlark.String:
		return string(v), nil

	default:
		return nil, errors.New("starlarkToGo: Unhandled type " + reflect.TypeOf(v).String())
	}
}
