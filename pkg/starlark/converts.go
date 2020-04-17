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
