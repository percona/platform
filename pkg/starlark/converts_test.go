package starlark

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoToStarlark(t *testing.T) {
	data := make(map[string]interface{})
	data["bool"] = true
	data["string"] = "Test string"
	data["int64"] = int64(-9045646465464654500)
	data["uint64"] = uint64(18446744073709551615)
	data["float64"] = float64(5.5)
	data["timestamp"] = time.Now().UnixNano()
	data["bytes"] = make([]byte, len(data["string"].(string)))
	data["bytes"] = []byte(data["string"].(string))
	data["none"] = nil

	for k, v := range data {
		k := k
		v := v
		t.Run(k, func(t *testing.T) {
			sv, errIn := goToStarlark(v)
			require.NoError(t, errIn)
			gv, errOut := starlarkToGo(sv)
			require.NoError(t, errOut)

			var res interface{}
			switch v.(type) {
			case []byte:
				res = []byte(gv.(string))
			default:
				res = gv
			}

			assert.Equal(t, v, res, "not equal ("+k+")")
		})
	}
}
