// Package starlark is executor for starklark.
package starlark

import (
	"testing"

	"github.com/percona-platform/platform/pkg/check"
)

func TestRun(t *testing.T) {
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

		if res.Status != check.Success {
			t.Error(res.Message)
		}
	})

	t.Run("float only", func(t *testing.T) {
		var data []map[string]interface{}
		data = append(data, dataFloat)
		res, _ := Run("float", script, "test", data)

		if res.Status != check.Success {
			t.Error(res.Message)
		}
	})

	t.Run("string only", func(t *testing.T) {
		var data []map[string]interface{}
		data = append(data, dataStr)
		res, _ := Run("string", script, "test", data)

		if res.Status != check.Success {
			t.Error(res.Message)
		}
	})

	t.Run("mixed", func(t *testing.T) {
		var data []map[string]interface{}
		data = append(data, dataStr)
		data = append(data, dataFloat)
		data = append(data, dataInt)
		res, _ := Run("mixed", script, "test", data)

		if res.Status != check.Success {
			t.Error(res.Message)
		}
	})
}
