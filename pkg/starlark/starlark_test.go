package starlark

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/percona-platform/platform/pkg/check"
)

func TestRun(t *testing.T) {
	p, err := os.Getwd()
	if err != nil {
		t.Error()
	}

	script := filepath.Join(p, "starlark_script.star")

	data_int := make(map[string]interface{})
	data_int["item1"] = 5
	data_int["item2"] = 10

	data_float := make(map[string]interface{})
	data_float["item3"] = 5.444
	data_float["item4"] = 10.111

	data_str := make(map[string]interface{})
	data_str["item5"] = "B"
	data_str["item6"] = "A"

	t.Run("int only", func(t *testing.T) {
		var data []map[string]interface{}
		data = append(data, data_int)
		res := Run("int", script, "test", data)

		if res.Status != check.Success {
			t.Error(res.Message)
		}
	})

	t.Run("float only", func(t *testing.T) {
		var data []map[string]interface{}
		data = append(data, data_float)
		res := Run("float", script, "test", data)

		if res.Status != check.Success {
			t.Error(res.Message)
		}
	})

	t.Run("string only", func(t *testing.T) {
		var data []map[string]interface{}
		data = append(data, data_str)
		res := Run("float", script, "test", data)

		if res.Status != check.Success {
			t.Error(res.Message)
		}
	})

	t.Run("mixed", func(t *testing.T) {
		var data []map[string]interface{}
		data = append(data, data_str)
		data = append(data, data_float)
		data = append(data, data_int)
		res := Run("float", script, "test", data)

		if res.Status != check.Success {
			t.Error(res.Message)
		}
	})
}
