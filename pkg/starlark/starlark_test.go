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

    results = []
    for row in rows:
        name, value = row["Variable_name"], row["Value"]
        expected = vars.get(name)
        if expected and expected != value:
            results.append({
                      "summary": "expected %s to be %s, got %s" % (name, expected, value),
                      "description": "description text",
                      "severity": "warning",
            })

    return results, ""
	`) + "\n"

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		input := []map[string]interface{}{
			{"Variable_name": "have_openssl", "Value": "YES"},
			{"Variable_name": "have_ssl", "Value": "YES"},
		}

		res, err := Run(t.Name(), script, input)
		require.NoError(t, err)
		assert.Empty(t, res)
	})

	t.Run("Single check result", func(t *testing.T) {
		t.Parallel()

		input := []map[string]interface{}{
			{"Variable_name": "have_ssl", "Value": "YES"},
			{"Variable_name": "have_openssl", "Value": "NO"},
		}

		res, err := Run(t.Name(), script, input)
		require.NoError(t, err)
		expected := []check.Result{
			{
				Severity:    check.Warning,
				Description: "description text",
				Summary:     "expected have_openssl to be YES, got NO",
			},
		}
		assert.Equal(t, expected, res)
	})

	t.Run("Multiple check results", func(t *testing.T) {
		t.Parallel()

		input := []map[string]interface{}{
			{"Variable_name": "have_ssl", "Value": "NO"},
			{"Variable_name": "have_openssl", "Value": "NO"},
		}

		res, err := Run(t.Name(), script, input)
		require.NoError(t, err)
		expected := []check.Result{
			{
				Severity:    check.Warning,
				Description: "description text",
				Summary:     "expected have_ssl to be YES, got NO",
			},
			{
				Severity:    check.Warning,
				Description: "description text",
				Summary:     "expected have_openssl to be YES, got NO",
			},
		}
		assert.Equal(t, expected, res)
	})
}
