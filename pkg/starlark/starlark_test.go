package starlark

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/percona-platform/platform/pkg/check"
)

func TestRunValidScript(t *testing.T) {
	t.Parallel()

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
        print(name, actual, expected)
        if expected and expected != actual:
            return {"error": "expected %s to be %s, got %s" % (name, expected, actual)}

	return {}
	`) + "\n"

	env, err := NewEnv(t.Name(), script)
	require.NoError(t, err)
	env.Print = func(msg string) { t.Log(msg) }

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		input := []map[string]interface{}{
			{"Variable_name": "have_ssl", "Value": "YES"},
			{"Variable_name": "have_openssl", "Value": "YES"},
		}

		res, err := env.Run(input)
		require.NoError(t, err)
		expected := &check.Result{
			Status:  "status",
			Message: "message",
		}
		assert.Equal(t, expected, res)
	})

	t.Run("Fail", func(t *testing.T) {
		t.Parallel()

		input := []map[string]interface{}{
			{"Variable_name": "have_ssl", "Value": "NO"},
			{"Variable_name": "have_openssl", "Value": "NO"},
		}

		res, err := env.Run(input)
		require.NoError(t, err)
		expected := &check.Result{
			Status:  "status",
			Message: "message",
		}
		assert.Equal(t, expected, res)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		res, err := env.Run(nil)
		assert.EqualError(t, err, "unhandled result type starlark.NoneType")
		assert.Nil(t, res)
	})
}
