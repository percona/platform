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
	case []byte:
		return starlark.String(v), nil
	case string:
		return starlark.String(v), nil
	case time.Time:
		return starlark.MakeInt64(v.UnixNano()), nil
	case []interface{}:
		return goToStarlarkList(v)
	case map[string]interface{}:
		return goToStarlarkDict(v)
	case map[bool]struct{}:
		return goToStarlarkSetBool(v)
	case map[int64]struct{}:
		return goToStarlarkSetInt(v)
	case map[uint64]struct{}:
		return goToStarlarkSetUint(v)
	case map[float64]struct{}:
		return goToStarlarkSetFloat(v)
	case map[string]struct{}:
		return goToStarlarkSetString(v)
	default:
		return nil, errors.Errorf("unhandled type %[1]T (%[1]v)", v)
	}
}

func goToStarlarkList(v []interface{}) (*starlark.List, error) {
	l := make([]starlark.Value, len(v))
	for i, o := range v {
		sv, err := goToStarlark(o)
		if err != nil {
			return nil, err
		}
		l[i] = sv
	}
	return starlark.NewList(l), nil
}

func goToStarlarkDict(v map[string]interface{}) (*starlark.Dict, error) {
	sd := starlark.NewDict(len(v))
	for k, o := range v {
		sv, err := goToStarlark(o)
		if err != nil {
			return nil, err
		}
		if err := sd.SetKey(starlark.String(k), sv); err != nil {
			return nil, errors.Wrap(err, "failed to set key in dict")
		}
	}
	return sd, nil
}

func goToStarlarkSetBool(v map[bool]struct{}) (*starlark.Set, error) {
	ss := starlark.NewSet(len(v))
	for k := range v {
		if err := ss.Insert(starlark.Bool(k)); err != nil {
			return nil, errors.Wrap(err, "failed to insert into set")
		}
	}
	return ss, nil
}

func goToStarlarkSetInt(v map[int64]struct{}) (*starlark.Set, error) {
	ss := starlark.NewSet(len(v))
	for k := range v {
		if err := ss.Insert(starlark.MakeInt64(k)); err != nil {
			return nil, errors.Wrap(err, "failed to insert into set")
		}
	}
	return ss, nil
}

func goToStarlarkSetUint(v map[uint64]struct{}) (*starlark.Set, error) {
	ss := starlark.NewSet(len(v))
	for k := range v {
		if err := ss.Insert(starlark.MakeUint64(k)); err != nil {
			return nil, errors.Wrap(err, "failed to insert into set")
		}
	}
	return ss, nil
}

func goToStarlarkSetFloat(v map[float64]struct{}) (*starlark.Set, error) {
	ss := starlark.NewSet(len(v))
	for k := range v {
		if err := ss.Insert(starlark.Float(k)); err != nil {
			return nil, errors.Wrap(err, "failed to insert into set")
		}
	}
	return ss, nil
}

func goToStarlarkSetString(v map[string]struct{}) (*starlark.Set, error) {
	ss := starlark.NewSet(len(v))
	for k := range v {
		if err := ss.Insert(starlark.String(k)); err != nil {
			return nil, errors.Wrap(err, "failed to insert into set")
		}
	}
	return ss, nil
}

// starlarkToGo converts Starlark value to Go value.
// Supported types:
//  * NoneType -> nil;
//  * bool -> bool;
//  * int -> int64/uint64;
//  * float -> float64;
//  * string -> string;
//  * tuple -> []interface{}
//  * list -> []interface{}
//  * dict (with string keys) -> map[string]interface{}
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
		return nil, errors.Errorf("interger value %s is too big", v)

	case starlark.Float:
		return float64(v), nil

	case starlark.String:
		return string(v), nil

	case starlark.Tuple:
		res := make([]interface{}, len(v))
		for i, el := range v {
			gv, err := starlarkToGo(el)
			if err != nil {
				return nil, err
			}
			res[i] = gv
		}
		return res, nil

	case *starlark.List:
		res := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			gv, err := starlarkToGo(v.Index(i))
			if err != nil {
				return nil, err
			}
			res[i] = gv
		}
		return res, nil

	case *starlark.Dict:
		res := make(map[string]interface{}, v.Len())
		for _, tu := range v.Items() {
			k, v := tu[0], tu[1]
			ks, ok := k.(starlark.String)
			if !ok {
				return nil, errors.Errorf("unhandled dict key type %[1]T (%[1]v)", k)
			}
			gv, err := starlarkToGo(v)
			if err != nil {
				return nil, err
			}
			res[string(ks)] = gv
		}
		return res, nil

	default:
		return nil, errors.Errorf("unhandled type %T", v)
	}
}
