package starlark

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoToStarlark(t *testing.T) {
	number := 500
	data := make(map[string]interface{})
	data["bool"] = true
	data["string"] = "Test string"
	data["int"] = int(number)
	data["int8"] = int8(number)
	data["int16"] = int16(number)
	data["int32"] = int32(number)
	data["int64"] = int64(number)

	for k, v := range data {
		t.Run(k, func(t *testing.T) {
			sv, errIn := goToStarlark(v)
			gv, errOut := starlarkToGo(sv)
			require.NoError(t, errIn)
			require.NoError(t, errOut)

			var res interface{}
			switch v := v.(type) {
			case int8:
				res = int(v)
			case int16:
				res = int(v)
			case int32:
				res = int(v)
			case int64:
				res = int(v)
			default:
				res = v
			}

			assert.Equal(t, res, gv, "not equal ("+k+")")
		})
	}
}
