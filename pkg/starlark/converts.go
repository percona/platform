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
		return starlark.Tuple(l), nil
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
	case map[int64]struct{}:
		ss := starlark.NewSet(len(v))
		for k := range v {
			err := ss.Insert(starlark.MakeInt64(k))
			if err != nil {
				return nil, errors.Wrap(err, "goToStarlark")
			}
		}
		return ss, nil
	case map[uint64]struct{}:
		ss := starlark.NewSet(len(v))
		for k := range v {
			err := ss.Insert(starlark.MakeUint64(k))
			if err != nil {
				return nil, errors.Wrap(err, "goToStarlark")
			}
		}
		return ss, nil
	case map[float64]struct{}:
		ss := starlark.NewSet(len(v))
		for k := range v {
			err := ss.Insert(starlark.Float(k))
			if err != nil {
				return nil, errors.Wrap(err, "goToStarlark")
			}
		}
		return ss, nil
	case map[string]struct{}:
		ss := starlark.NewSet(len(v))
		for k := range v {
			err := ss.Insert(starlark.String(k))
			if err != nil {
				return nil, errors.Wrap(err, "goToStarlark")
			}
		}
		return ss, nil
	default:
		return nil, errors.Errorf("goToStarlark: unhandled type %T", v)
	}
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
	case *starlark.Set:
		return starlarkSetToGo(v)
	default:
		return nil, errors.Errorf("starlarkToGo: unhandled type %T", v)
	}
}

func starlarkSetToGo(v *starlark.Set) (interface{}, error) {
	iter := v.Iterate()
	var x starlark.Value
	var tp string
	var notValid bool
	for iter.Next(&x) {
		if tp != "" && tp != x.Type() {
			notValid = true
			iter.Done()
		}
		tp = x.Type()
	}
	if notValid {
		return nil, errors.New("starlarkSetToGo: More types in starlark.Set")
	}
	iter = v.Iterate()
	defer iter.Done()

	switch tp {
	case "int":
		vl := x.(starlark.Int)
		if _, ok := vl.Int64(); ok {
			res := make(map[int64]struct{}, v.Len())
			for iter.Next(&x) {
				nv, err := starlarkToGo(x.(starlark.Int))
				if err != nil {
					return nil, err
				}
				res[nv.(int64)] = struct{}{}
			}
			return res, nil
		}
		if _, ok := vl.Uint64(); ok {
			res := make(map[uint64]struct{}, v.Len())
			for iter.Next(&x) {
				nv, err := starlarkToGo(x.(starlark.Int))
				if err != nil {
					return nil, err
				}
				res[nv.(uint64)] = struct{}{}
			}
			return res, nil
		}
		return nil, errors.Errorf("starlarkSetToGo: unhandled type %s", tp)
	case "float":
		res := make(map[float64]struct{}, v.Len())
		for iter.Next(&x) {
			nv, err := starlarkToGo(x.(starlark.Float))
			if err != nil {
				return nil, err
			}
			res[nv.(float64)] = struct{}{}
		}
		return res, nil
	case "string":
		res := make(map[string]struct{}, v.Len())
		for iter.Next(&x) {
			res[x.(starlark.String).GoString()] = struct{}{}
		}
		return res, nil
	default:
		return nil, errors.Errorf("starlarkSetToGo: unhandled type %s", tp)
	}
}
