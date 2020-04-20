// Package starlark is executor for starklark.
package starlark

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/percona-platform/platform/pkg/check"
)

func TestRun(t *testing.T) {
	assert := assert.New(t)
	script := "starlark_script.py"

	dataInt := make(map[string]interface{})
	dataInt["item1"] = 5
	dataInt["item2"] = 10

	dataFloat := make(map[string]interface{})
	dataFloat["item3"] = 5.444
	dataFloat["item4"] = 10.111

	dataStr := make(map[string]interface{})
	dataStr["item5"] = "B"
	dataStr["item6"] = "A"

	t.Run("int only", func(t *testing.T) {
		var data []map[string]interface{}
		data = append(data, dataInt)
		res, _ := Run("int", script, "test", data)

		assert.Equal(res.Status, check.Success, res.Message)
	})

	t.Run("float only", func(t *testing.T) {
		var data []map[string]interface{}
		data = append(data, dataFloat)
		res, _ := Run("float", script, "test", data)

		assert.Equal(res.Status, check.Success, res.Message)
	})

	t.Run("string only", func(t *testing.T) {
		var data []map[string]interface{}
		data = append(data, dataStr)
		res, _ := Run("string", script, "test", data)

		assert.Equal(res.Status, check.Success, res.Message)
	})

	t.Run("mixed", func(t *testing.T) {
		var data []map[string]interface{}
		data = append(data, dataStr)
		data = append(data, dataFloat)
		data = append(data, dataInt)
		res, _ := Run("mixed", script, "test", data)

		assert.Equal(res.Status, check.Success, res.Message)
	})

	script = "starlark_script2.py"
	t.Run("check", func(t *testing.T) {
		dataCheck := make(map[string]interface{})
		dataCheck["Variable_name"] = "have_ssl"
		dataCheck["Value"] = "YES"

		var data []map[string]interface{}
		data = append(data, dataCheck)
		res, _ := Run("check", script, "check", data)

		assert.Equal(res.Status, check.Success, res.Message)
	})
}
