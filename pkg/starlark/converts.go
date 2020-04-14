package starlark

import (
	"fmt"

	"go.starlark.net/starlark"
)

func goToStarlark(v interface{}) starlark.Value {
	switch v := v.(type) {
	case int:
		return starlark.MakeInt(v)
	case string:
		return starlark.String(v)
	case float32:
		return starlark.Float(v)
	case float64:
		return starlark.Float(v)
	default:
		panic(fmt.Sprintf("goToStarlark: unhandled %#[1]v (%[1]T)", v))
	}
}

func starlarkToGo(v starlark.Value) interface{} {
	switch v := v.(type) {
	case starlark.Int:
		if i, ok := v.Int64(); ok {
			return int(i)
		}
		if u, ok := v.Uint64(); ok {
			return uint(u)
		}
		panic(fmt.Sprintf("starlarkToGo: unhandled starlark.Int %s", v.String()))
	case starlark.Float:
		return v
	case starlark.String:
		return string(v)
	default:
		panic(fmt.Sprintf("starlarkToGo: unhandled %#[1]v (%[1]T)", v))
	}
}
