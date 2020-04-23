package starlark

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	script := "starlark_script.py"

	dataInt := make(map[string]interface{})
	dataInt["item1"] = int64(5)
	dataInt["item2"] = int64(10)

	dataFloat := make(map[string]interface{})
	dataFloat["item3"] = float64(5.444)
	dataFloat["item4"] = float64(10.111)

	dataStr := make(map[string]interface{})
	dataStr["item5"] = "B"
	dataStr["item6"] = "A"

	t.Run("int only", func(t *testing.T) {
		t.Parallel()

		var data []map[string]interface{}
		data = append(data, dataInt)
		_, err := run("int", script, "test", data)

		require.NoError(t, err)
	})

	t.Run("float only", func(t *testing.T) {
		t.Parallel()

		var data []map[string]interface{}
		data = append(data, dataFloat)
		_, err := run("float", script, "test", data)

		require.NoError(t, err)
	})

	t.Run("string only", func(t *testing.T) {
		t.Parallel()

		var data []map[string]interface{}
		data = append(data, dataStr)
		_, err := run("string", script, "test", data)

		require.NoError(t, err)
	})

	t.Run("mixed", func(t *testing.T) {
		t.Parallel()

		var data []map[string]interface{}
		data = append(data, dataStr)
		data = append(data, dataFloat)
		data = append(data, dataInt)
		_, err := run("mixed", script, "test", data)

		require.NoError(t, err)
	})

	t.Run("check", func(t *testing.T) {
		t.Parallel()

		dataCheck := make(map[string]interface{})
		dataCheck["Variable_name"] = "have_ssl"
		dataCheck["Value"] = "YES"

		var data []map[string]interface{}
		data = append(data, dataCheck)
		_, err := run("check", script, "check", data)

		require.NoError(t, err)
	})
}
