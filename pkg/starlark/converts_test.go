package starlark

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoToStarlark(t *testing.T) {
	data := map[string]interface{}{
		"none":       nil,
		"bool":       true,
		"int64":      int64(-9045646465464654500),
		"uint64":     uint64(18446744073709551615),
		"float64":    float64(5.5),
		"bytes":      []byte("Test"),
		"string":     "Test",
		"timestamp":  time.Now().Truncate(0),
		"slice":      []interface{}{int64(500), "Test", float64(30.555555555555)},
		"map":        map[string]interface{}{"ka": "a", "kb": "b", "kc": "c", "kd": "d"},
		"boolset":    map[bool]struct{}{true: {}, false: {}},
		"int64set":   map[int64]struct{}{50: {}, 20: {}},
		"float64set": map[float64]struct{}{50.55555: {}, 10.2456789: {}},
		"stringset":  map[string]struct{}{"test": {}, "test2": {}},
	}

	for k, v := range data {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			sv, errIn := goToStarlark(v)
			require.NoError(t, errIn)
			gv, errOut := starlarkToGo(sv)
			require.NoError(t, errOut)

			var res interface{}
			switch v.(type) {
			case time.Time:
				res = time.Unix(0, gv.(int64))
			case []byte:
				res = []byte(gv.(string))
			default:
				res = gv
			}

			assert.Equal(t, v, res, "not equal ("+k+")")
		})
	}
}

func TestStarlarkToGo(t *testing.T) {
	t.Parallel()

	input := []interface{}{
		string("Test"),
		int64(-9045646465464654500),
		uint64(18446744073709551615),
		float64(5.555555555555),
	}

	rows, errIn := goToStarlark(input)
	require.NoError(t, errIn)

	gr, errOut := starlarkToGo(rows)
	require.NoError(t, errOut)

	assert.Equal(t, input, gr, "not equal")
}
