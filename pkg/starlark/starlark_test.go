package starlark

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/percona-platform/platform/pkg/check"
)

func TestRun(t *testing.T) {
	script := strings.TrimSpace(`
def check(rows):
    vars = {
        "have_ssl":     "YES",
        "have_openssl": "YES",
    }

    for row in rows:
        name = row["Variable_name"]
        actual = row["Value"]
        expected = vars.get(name)
        if expected and expected != actual:
            return {"error": "expected %s to be %s, got %s" % (name, expected, actual)}

	return {}
	`) + "\n"

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		input := []map[string]interface{}{
			{"Variable_name": "have_ssl", "Value": "YES"},
			{"Variable_name": "have_openssl", "Value": "YES"},
		}

		res, err := Run(t.Name(), script, input)
		require.NoError(t, err)
		expected := &check.Result{
			Severity: check.Info,
			Summary:  "summary",
		}
		assert.Equal(t, expected, res)
	})

	t.Run("Fail", func(t *testing.T) {
		t.Parallel()

		input := []map[string]interface{}{
			{"Variable_name": "have_ssl", "Value": "NO"},
			{"Variable_name": "have_openssl", "Value": "NO"},
		}

		res, err := Run(t.Name(), script, input)
		require.NoError(t, err)
		expected := &check.Result{
			Severity: check.Info,
			Summary:  "summary",
		}
		assert.Equal(t, expected, res)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		res, err := Run(t.Name(), script, nil)
		assert.EqualError(t, err, "unhandled result type starlark.NoneType")
		assert.Nil(t, res)
	})
}
