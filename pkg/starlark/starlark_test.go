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

		res, err := env.Run(t.Name(), input)
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

		res, err := env.Run(t.Name(), input)
		require.NoError(t, err)
		expected := &check.Result{
			Status:  "status",
			Message: "message",
		}
		assert.Equal(t, expected, res)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		res, err := env.Run(t.Name(), nil)
		assert.EqualError(t, err, "unhandled result type starlark.NoneType")
		assert.Nil(t, res)
	})
}

func TestRunInvalidScript(t *testing.T) {
	t.Parallel()

	t.Run("Parse", func(t *testing.T) {
		t.Parallel()

		env, err := NewEnv(t.Name(), `foo`)
		assert.Nil(t, env)

		expected := "failed to parse script: TestRunInvalidScript/Parse:1:1: undefined: foo"
		assert.EqualError(t, err, expected)
	})

	t.Run("Undefined", func(t *testing.T) {
		t.Parallel()

		env, err := NewEnv(t.Name(), `def foo(): pass`)
		require.NoError(t, err)

		res, err := env.run("id", "bar", nil)
		assert.Nil(t, res)

		expected := "id: function bar is not defined"
		assert.EqualError(t, err, expected)

		expected = "id: function bar is not defined"
		assert.Equal(t, expected, fmt.Sprintf("%+v", err))
	})

	t.Run("Execute", func(t *testing.T) {
		t.Parallel()

		env, err := NewEnv(t.Name(), `def foo(): 0/0`)
		require.NoError(t, err)

		res, err := env.run("id", "foo", nil)
		assert.Nil(t, res)

		expected := strings.TrimSpace(`
id: failed to execute function foo
Traceback (most recent call last):
  TestRunInvalidScript/Execute:1:13: in foo
: real division by zero
		`)
		assert.EqualError(t, err, expected)

		expected = strings.TrimSpace(`
id: failed to execute function foo
Traceback (most recent call last):
  TestRunInvalidScript/Execute:1:13: in foo
: real division by zero
		`)
		assert.Equal(t, expected, fmt.Sprintf("%+v", err))
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

	env, err := NewEnv(t.Name(), script)
	require.NoError(t, err)

	var buf bytes.Buffer
	env.Print = func(msg string) { _, _ = buf.WriteString(msg) }

	res, err := env.run("id", "test1", nil)
	require.NoError(t, err)
	assert.Equal(t, starlark.None, res)

	expected := strings.TrimSpace(`
id TestPrint:8:6 <toplevel>: hello from main
id TestPrint:5:10 test1: hello from test1
id TestPrint:2:10 test2: hello from test2
	`) + "\n"
	assert.Equal(t, expected, buf.String())
}
