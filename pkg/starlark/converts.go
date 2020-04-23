package starlark

import (
	"time"

	"github.com/pkg/errors"
	"go.starlark.net/starlark"
)

func goToStarlark(v interface{}) (starlark.Value, error) {
	switch v := v.(type) {
	case nil:
		return starlark.None, nil
	case bool:
		return starlark.Bool(v), nil
	case int64:
		return starlark.MakeInt64(v), nil
	case uint64:
		return starlark.MakeUint64(v), nil
	case float64:
		return starlark.Float(v), nil
	case string:
		return starlark.String(v), nil
	case []byte:
		return starlark.String(v), nil
	case time.Time:
		return starlark.MakeInt64(v.UnixNano()), nil
	case []interface{}:
		var l []starlark.Value
		for _, o := range v {
			sv, err := goToStarlark(o)
			if err != nil {
				return nil, err
			}
			l = append(l, sv)
		}
		return starlark.NewList(l), nil
	case map[string]interface{}:
		sd := starlark.NewDict(len(v))
		for k, o := range v {
			sv, err := goToStarlark(o)
			if err != nil {
				return nil, err
			}
			if err := sd.SetKey(starlark.String(k), sv); err != nil {
				return nil, errors.Wrap(err, "goToStarlark")
			}
		}
		return sd, nil
	}

	return nil, errors.Errorf("goToStarlark: unhandled type %T", v)
}

func starlarkToGo(v starlark.Value) (interface{}, error) {
	switch v := v.(type) {
	case starlark.NoneType:
		return nil, nil
	case starlark.Bool:
		return bool(v), nil
	case starlark.Int:
		if i, ok := v.Int64(); ok {
			return i, nil
		}
		if u, ok := v.Uint64(); ok {
			return u, nil
		}
		return nil, errors.Errorf("starlarkToGo: unhandled type %T", v)
	case starlark.Float:
		return float64(v), nil
	case starlark.String:
		return string(v), nil
	case *starlark.List:
		var res []interface{}
		for i := 0; i < v.Len(); i++ {
			gv, err := starlarkToGo(v.Index(i))
			if err != nil {
				return nil, err
			}
			res = append(res, gv)
		}
		return res, nil
	case *starlark.Dict:
		res := make(map[string]interface{}, v.Len())
		for _, tu := range v.Items() {
			var err error
			k := tu[0].(starlark.String).GoString()
			res[k], err = starlarkToGo(tu[1])
			if err != nil {
				return nil, err
			}
		}
		return res, nil
	case starlark.Tuple:
		var res []interface{}
		for _, o := range v {
			no, err := starlarkToGo(o)
			if err != nil {
				return nil, err
			}
			res = append(res, no)
		}
		return res, nil
	}

	return nil, errors.Errorf("starlarkToGo: unhandled type %T", v)
}
