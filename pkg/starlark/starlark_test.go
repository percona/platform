package starlark

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.starlark.net/starlark"

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

    results = []
    for row in rows:
        name, actual = row["Variable_name"], row["Value"]
        expected = vars.get(name)
        print(name, actual, expected)
        if expected and expected != actual:
            results.append({
                      "summary": "expected %s to be %s, got %s" % (name, expected, actual),
                      "description": "description text",
                      "severity": "warning",
            })

    return results, ""
	`) + "\n"

	addToFuzzCorpus(t.Name(), script, nil)
	env, err := NewEnv(t.Name(), script)
	require.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		input := []map[string]interface{}{
			{"Variable_name": "have_openssl", "Value": "YES"},
			{"Variable_name": "have_ssl", "Value": "YES"},
		}

		addToFuzzCorpus(t.Name(), script, input)
		res, err := env.Run(t.Name(), input, t.Log)
		require.NoError(t, err)
		assert.Empty(t, res)
	})

	t.Run("Single check result", func(t *testing.T) {
		t.Parallel()

		input := []map[string]interface{}{
			{"Variable_name": "have_ssl", "Value": "YES"},
			{"Variable_name": "have_openssl", "Value": "NO"},
		}

		addToFuzzCorpus(t.Name(), script, input)
		res, err := env.Run(t.Name(), input, t.Log)
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

		addToFuzzCorpus(t.Name(), script, input)
		res, err := env.Run(t.Name(), input, t.Log)
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

func TestRunInvalidScript(t *testing.T) {
	t.Parallel()

	t.Run("Parse", func(t *testing.T) {
		t.Parallel()

		script := `def foo(): parse_version("2.6.0")`
		addToFuzzCorpus(t.Name(), script, nil)
		env, err := NewEnv(t.Name(), script)
		assert.Nil(t, env)

		expected := "failed to parse script: TestRunInvalidScript/Parse:1:12: undefined: parse_version"
		assert.EqualError(t, err, expected)
	})

	t.Run("Init", func(t *testing.T) {
		t.Parallel()

		script := `""[1]`
		addToFuzzCorpus(t.Name(), script, nil)
		env, err := NewEnv(t.Name(), script)
		require.NoError(t, err)

		res, err := env.run("bar", nil, "id", t.Log)
		assert.Nil(t, res)

		expected := strings.TrimSpace(`
[id] failed to init script: index 1 out of range: empty string
Traceback (most recent call last):
  TestRunInvalidScript/Init:1:3: in <toplevel>
		`) + "\n"
		assert.EqualError(t, err, expected)
	})

	t.Run("Undefined", func(t *testing.T) {
		t.Parallel()

		script := `def foo(): pass`
		addToFuzzCorpus(t.Name(), script, nil)
		env, err := NewEnv(t.Name(), script)
		require.NoError(t, err)

		res, err := env.run("bar", nil, "id", t.Log)
		assert.Nil(t, res)

		expected := "[id] function bar is not defined"
		assert.EqualError(t, err, expected)
	})

	t.Run("Execute", func(t *testing.T) {
		t.Parallel()

		script := `def foo(): 0/0`
		addToFuzzCorpus(t.Name(), script, nil)
		env, err := NewEnv(t.Name(), script)
		require.NoError(t, err)

		res, err := env.run("foo", nil, "id", t.Log)
		assert.Nil(t, res)

		expected := strings.TrimSpace(`
[id] failed to execute function foo: real division by zero
Traceback (most recent call last):
  TestRunInvalidScript/Execute:1:13: in foo
		`) + "\n"
		assert.EqualError(t, err, expected)
	})

	t.Run("Hang", func(t *testing.T) {
		t.Parallel()

		script := `[7]*714748364`
		addToFuzzCorpus(t.Name(), script, nil)
		env, err := NewEnv(t.Name(), script)
		require.NoError(t, err)

		_ = env
		t.Skip("https://jira.percona.com/browse/SAAS-63")

		_, err = env.run("foo", nil, "id", t.Log)
		assert.EqualError(t, err, "context timeout or something")
	})
}

func TestPrint(t *testing.T) {
	t.Parallel()

	script := strings.TrimSpace(`
def test2():
    print("hello from test2")

def test1():
    print("hello from test1")
    test2()

print("hello from main")
	`) + "\n"

	addToFuzzCorpus(t.Name(), script, nil)
	env, err := NewEnv(t.Name(), script)
	require.NoError(t, err)

	var buf bytes.Buffer
	print := func(args ...interface{}) {
		_, _ = buf.WriteString(fmt.Sprintln(args...))
	}

	res, err := env.run("test1", nil, "id", print)
	require.NoError(t, err)
	assert.Equal(t, starlark.None, res)

	expected := strings.TrimSpace(`
[id] TestPrint:8:6: in <toplevel>: hello from main
[id] TestPrint:5:10: in test1: hello from test1
[id] TestPrint:2:10: in test2: hello from test2
	`) + "\n"
	assert.Equal(t, expected, buf.String())
}
