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
		t.Parallel()

		for _, v := range []interface{}{
			nil,
			true,
			int64(-9045646465464654500),
			uint64(18446744073709551615),
			float64(5.5),
			"Test",
			[]interface{}{int64(500), "Test", float64(30.555555555555)},
			map[string]interface{}{"ka": "a", "kb": "b", "kc": "c", "kd": "d"},
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

	t.Run("goToStarlark", func(t *testing.T) {
		t.Parallel()

		type pair struct {
			gv interface{}
			sv starlark.Value
		}

		for _, p := range []pair{
			{[]byte("Test"), starlark.String("Test")},
		} {
			v, expected := p.gv, p.sv
			t.Run(fmt.Sprint(v), func(t *testing.T) {
				t.Parallel()

				sv, err := goToStarlark(v)
				require.NoError(t, err)
				assert.Equal(t, expected, sv)
			})
		}

		// special case
		sv, err := goToStarlark(time.Date(2020, 4, 28, 13, 48, 42, 0, time.UTC))
		require.NoError(t, err)
		expected := starlark.MakeInt64(1588081722000000000).BigInt()
		actual := sv.(starlark.Int).BigInt()
		assert.Equal(t, expected, actual)
	})

	t.Run("starlarkToGo", func(t *testing.T) {
		t.Parallel()

		type pair struct {
			sv starlark.Value
			gv interface{}
		}

		for _, p := range []pair{
			{starlark.MakeInt64(1588081722000000000), int64(1588081722000000000)},
			{starlark.String("Test"), "Test"},
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

		t.Run("intDict", func(t *testing.T) {
			dict := starlark.NewDict(1)
			err := dict.SetKey(starlark.MakeInt(1), starlark.MakeInt(2))
			require.NoError(t, err)
			_, err = starlarkToGo(dict)
			assert.EqualError(t, err, "unhandled dict key type starlark.Int (1)")
		})
	})
}
