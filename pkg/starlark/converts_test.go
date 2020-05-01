package starlark

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.starlark.net/starlark"
)

func TestConvert(t *testing.T) {
	t.Parallel()

	t.Run("BothWays", func(t *testing.T) {
		for _, v := range []interface{}{
			nil,
			true,
			int64(-9045646465464654500),
			uint64(18446744073709551615),
			float64(5.5),
			"Test",
			[]interface{}{int64(500), "Test", float64(30.555555555555)},
			map[string]interface{}{"ka": "a", "kb": "b", "kc": "c", "kd": "d"},
			map[bool]struct{}{true: {}, false: {}},
			map[int64]struct{}{50: {}, 20: {}},
			map[float64]struct{}{50.55555: {}, 10.2456789: {}},
			map[string]struct{}{"test": {}, "test2": {}},
		} {
			v := v
			t.Run(fmt.Sprint(v), func(t *testing.T) {
				t.Parallel()

				sv, err := goToStarlark(v)
				require.NoError(t, err)
				gv, err := starlarkToGo(sv)
				require.NoError(t, err)
				assert.Equal(t, v, gv, "sv = %#[1]v %[1]T", sv)
			})
		}
	})

	uint64Set := starlark.NewSet(0)
	require.NoError(t, uint64Set.Insert(starlark.MakeUint64(50)))
	require.NoError(t, uint64Set.Insert(starlark.MakeUint64(20)))

	t.Run("goToStarlark", func(t *testing.T) {
		type pair struct {
			gv interface{}
			sv starlark.Value
		}

		for _, p := range []pair{
			{time.Date(2020, 4, 28, 13, 48, 42, 0, time.UTC), starlark.MakeInt64(1588081722000000000)},
			{[]byte("Test"), starlark.String("Test")},
			{map[uint64]struct{}{50: {}, 20: {}}, uint64Set},
		} {
			v, expected := p.gv, p.sv
			t.Run(fmt.Sprint(v), func(t *testing.T) {
				t.Parallel()

				sv, err := goToStarlark(v)
				require.NoError(t, err)
				assert.Equal(t, expected, sv)
			})
		}
	})

	t.Run("starlarkToGo", func(t *testing.T) {
		type pair struct {
			sv starlark.Value
			gv interface{}
		}

		for _, p := range []pair{
			{starlark.MakeInt64(1588081722000000000), int64(1588081722000000000)},
			{starlark.String("Test"), "Test"},
			{uint64Set, map[int64]struct{}{50: {}, 20: {}}},
			{starlark.Tuple{starlark.MakeInt(50), starlark.MakeInt(20)}, []interface{}{int64(50), int64(20)}},
		} {
			v, expected := p.sv, p.gv
			t.Run(fmt.Sprint(v), func(t *testing.T) {
				t.Parallel()

				gv, err := starlarkToGo(v)
				require.NoError(t, err)
				assert.Equal(t, expected, gv)
			})
		}
	})
}
