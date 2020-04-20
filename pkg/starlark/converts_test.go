package starlark

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoToStarlark(t *testing.T) {
	data := make(map[string]interface{})
	data["bool"] = true
	data["string"] = "Test string"
	data["int64"] = int64(-500)
	data["uint64"] = uint64(500)
	data["float64"] = float64(5.5)

	for k, v := range data {
		k := k
		v := v
		t.Run(k, func(t *testing.T) {
			sv, errIn := goToStarlark(v)
			gv, errOut := starlarkToGo(sv)
			require.NoError(t, errIn)
			require.NoError(t, errOut)

			var res interface{}
			switch v.(type) {
			case uint64:
				res = uint64(gv.(int64))
			default:
				res = gv
			}

			assert.Equal(t, v, res, "not equal ("+k+")")
		})
	}
}
