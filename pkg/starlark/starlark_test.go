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
	"github.com/percona-platform/platform/pkg/common"
)

func TestRunValidScript(t *testing.T) {
	t.Parallel()

	script := strings.TrimSpace(`
def check(rows):
	return check_context(rows, {})

def check_context(rows, context):
    if not rows:
        return "no rows in result"

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
                "summary": "MySQL is not secured",
				"description": "expected {} to be {}, got {}".format(name, expected, actual),
				"read_more_url": "https://www.percona.com",
                "severity": "warning",
                "labels": {
                    name: actual,
                },
            })

    return results
	`) + "\n"

	addToFuzzCorpus(t.Name(), script, nil)
	env, err := NewEnv(t.Name(), script, nil)
	require.NoError(t, err)

	t.Run("NoResults", func(t *testing.T) {
		t.Parallel()

		data := []map[string]interface{}{
			{"Variable_name": "have_openssl", "Value": "YES"},
			{"Variable_name": "have_ssl", "Value": "YES"},
		}

		addToFuzzCorpus(t.Name(), script, data)
		input, err := PrepareInput(data)
		require.NoError(t, err)

		res, err := env.Run("id", input, nil, t.Log)
		require.NoError(t, err)
		assert.Empty(t, res)
	})

	t.Run("SingleResult", func(t *testing.T) {
		t.Parallel()

		data := []map[string]interface{}{
			{"Variable_name": "have_ssl", "Value": "YES"},
			{"Variable_name": "have_openssl", "Value": "NO"},
		}

		addToFuzzCorpus(t.Name(), script, data)
		input, err := PrepareInput(data)
		require.NoError(t, err)

		res, err := env.Run("id", input, nil, t.Log)
		require.NoError(t, err)
		expected := []check.Result{{
			Summary:     "MySQL is not secured",
			Description: "expected have_openssl to be YES, got NO",
			ReadMoreURL: "https://www.percona.com",
			Severity:    common.Warning,
			Labels:      map[string]string{"have_openssl": "NO"},
		}}
		assert.Equal(t, expected, res)
	})

	t.Run("MultipleResults", func(t *testing.T) {
		t.Parallel()

		data := []map[string]interface{}{
			{"Variable_name": "have_ssl", "Value": "NO"},
			{"Variable_name": "have_openssl", "Value": "NO"},
		}

		addToFuzzCorpus(t.Name(), script, data)
		input, err := PrepareInput(data)
		require.NoError(t, err)

		res, err := env.Run("id", input, nil, t.Log)
		require.NoError(t, err)
		expected := []check.Result{{
			Summary:     "MySQL is not secured",
			Description: "expected have_ssl to be YES, got NO",
			ReadMoreURL: "https://www.percona.com",
			Severity:    common.Warning,
			Labels:      map[string]string{"have_ssl": "NO"},
		}, {
			Summary:     "MySQL is not secured",
			Description: "expected have_openssl to be YES, got NO",
			ReadMoreURL: "https://www.percona.com",
			Severity:    common.Warning,
			Labels:      map[string]string{"have_openssl": "NO"},
		}}
		assert.Equal(t, expected, res)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		_, err := env.Run("id", nil, nil, t.Log)
		require.EqualError(t, err, "thread id: script returned error: no rows in result")
	})
}

func TestRunInvalidScript(t *testing.T) {
	t.Parallel()

	t.Run("Parse", func(t *testing.T) {
		t.Parallel()

		script := `def foo(): parse_version("2.6.0")`
		addToFuzzCorpus(t.Name(), script, nil)
		env, err := NewEnv(t.Name(), script, nil)
		assert.Nil(t, env)

		expected := `failed to parse script: TestRunInvalidScript/Parse:1:12: undefined: parse_version`
		assert.EqualError(t, err, expected)
	})

	t.Run("Init", func(t *testing.T) {
		t.Parallel()

		script := `""[1]`
		addToFuzzCorpus(t.Name(), script, nil)
		env, err := NewEnv(t.Name(), script, nil)
		require.NoError(t, err)

		res, err := env.run("bar", nil, "id", t.Log)
		assert.Nil(t, res)

		expected := strings.TrimSpace(`
thread id: failed to init script: index 1 out of range: empty string
Traceback (most recent call last):
  TestRunInvalidScript/Init:1:3: in <toplevel>
		`) + "\n"
		assert.EqualError(t, err, expected)
	})

	t.Run("Undefined", func(t *testing.T) {
		t.Parallel()

		script := `def foo(): pass`
		addToFuzzCorpus(t.Name(), script, nil)
		env, err := NewEnv(t.Name(), script, nil)
		require.NoError(t, err)

		res, err := env.run("bar", nil, "id", t.Log)
		assert.Nil(t, res)

		expected := `thread id: function bar is not defined`
		assert.EqualError(t, err, expected)
	})

	t.Run("Execute", func(t *testing.T) {
		t.Parallel()

		script := `def foo(): 0/0`
		addToFuzzCorpus(t.Name(), script, nil)
		env, err := NewEnv(t.Name(), script, nil)
		require.NoError(t, err)

		res, err := env.run("foo", nil, "id", t.Log)
		assert.Nil(t, res)

		expected := strings.TrimSpace(`
thread id: failed to execute function foo: floating-point division by zero
Traceback (most recent call last):
  TestRunInvalidScript/Execute:1:13: in foo
		`) + "\n"
		assert.EqualError(t, err, expected)
	})

	t.Run("Hang", func(t *testing.T) {
		t.Parallel()

		script := `def foo(): return [1] * (1 << 30-1)` // one less that maxAlloc in starlark
		addToFuzzCorpus(t.Name(), script, nil)
		env, err := NewEnv(t.Name(), script, nil)
		require.NoError(t, err)

		_ = env
		t.Skip("https://jira.percona.com/browse/SAAS-63")

		_, err = env.run("foo", nil, "id", t.Log)
		assert.EqualError(t, err, `context timeout or something`)
	})

	t.Run("InvalidOutputValue", func(t *testing.T) {
		t.Parallel()

		script := strings.TrimSpace(`
def check(rows):
	return check_context(rows, {})

def check_context(rows, context):
    return set()
		`) + "\n"

		addToFuzzCorpus(t.Name(), script, nil)
		env, err := NewEnv(t.Name(), script, nil)
		require.NoError(t, err)

		_, err = env.Run("id", nil, nil, t.Log)
		assert.EqualError(t, err, `thread id: failed to parse script output: unhandled type *starlark.Set`)
	})

	t.Run("InvalidOutputNotList", func(t *testing.T) {
		t.Parallel()

		script := strings.TrimSpace(`
def check(rows):
	return check_context(rows, {})

def check_context(rows, context):
    return {"summary": "foo"}
		`) + "\n"

		addToFuzzCorpus(t.Name(), script, nil)
		env, err := NewEnv(t.Name(), script, nil)
		require.NoError(t, err)

		_, err = env.Run("id", nil, nil, t.Log)
		assert.EqualError(t, err, `thread id: failed to parse script output: map[summary:foo] (map[string]interface {})`)
	})

	t.Run("InvalidOutputNotDict", func(t *testing.T) {
		t.Parallel()

		script := strings.TrimSpace(`
def check(rows):
	return check_context(rows, {})

def check_context(rows, context):
    return [1]
		`) + "\n"

		addToFuzzCorpus(t.Name(), script, nil)
		env, err := NewEnv(t.Name(), script, nil)
		require.NoError(t, err)

		_, err = env.Run("id", nil, nil, t.Log)
		assert.EqualError(t, err, `thread id: failed to parse script output: result 0 has wrong type: int64`)
	})

	t.Run("InvalidOutputNotString", func(t *testing.T) {
		t.Parallel()

		script := strings.TrimSpace(`
def check(rows):
	return check_context(rows, {})

def check_context(rows, context):
    return [{"summary": 1}]
		`) + "\n"

		addToFuzzCorpus(t.Name(), script, nil)
		env, err := NewEnv(t.Name(), script, nil)
		require.NoError(t, err)

		_, err = env.Run("id", nil, nil, t.Log)
		assert.EqualError(t, err, `thread id: failed to parse script output: "summary" has wrong type: int64 (1)`)
	})

	t.Run("InvalidResult", func(t *testing.T) {
		t.Parallel()

		script := strings.TrimSpace(`
def check(rows):
	return check_context(rows, {})

def check_context(rows, context):
    return [{}]
		`) + "\n"

		addToFuzzCorpus(t.Name(), script, nil)
		env, err := NewEnv(t.Name(), script, nil)
		require.NoError(t, err)

		_, err = env.Run("id", nil, nil, t.Log)
		assert.EqualError(t, err, `thread id: failed to parse script output: summary is empty`)
	})

	t.Run("InvalidLabels", func(t *testing.T) {
		t.Parallel()

		script := strings.TrimSpace(`
def check(rows):
	return check_context(rows, {})

def check_context(rows, context):
    return [{"labels": 1}]
		`) + "\n"

		addToFuzzCorpus(t.Name(), script, nil)
		env, err := NewEnv(t.Name(), script, nil)
		require.NoError(t, err)

		_, err = env.Run("id", nil, nil, t.Log)
		assert.EqualError(t, err, `thread id: failed to parse script output: labels field has wrong type: int64 (1)`)
	})

	t.Run("InvalidLabel", func(t *testing.T) {
		t.Parallel()

		script := strings.TrimSpace(`
def check(rows):
	return check_context(rows, {})

def check_context(rows, context):
    return [{"labels": {"foo": 1}}]
		`) + "\n"

		addToFuzzCorpus(t.Name(), script, nil)
		env, err := NewEnv(t.Name(), script, nil)
		require.NoError(t, err)

		_, err = env.Run("id", nil, nil, t.Log)
		assert.EqualError(t, err, `thread id: failed to parse script output: labels: "foo" has wrong type: int64 (1)`)
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
	env, err := NewEnv(t.Name(), script, nil)
	require.NoError(t, err)

	var buf bytes.Buffer
	printFunc := func(args ...interface{}) {
		_, _ = buf.WriteString(fmt.Sprintln(args...))
	}

	res, err := env.run("test1", nil, "id", printFunc)
	require.NoError(t, err)
	assert.Equal(t, starlark.None, res)

	expected := strings.TrimSpace(`
thread id: TestPrint:8:6: in <toplevel>: hello from main
thread id: TestPrint:5:10: in test1: hello from test1
thread id: TestPrint:2:10: in test2: hello from test2
	`) + "\n"
	assert.Equal(t, expected, buf.String())
}

func TestRegisterFunc(t *testing.T) {
	t.Parallel()

	pairs := func(args ...interface{}) (interface{}, error) {
		t.Logf("args = %#v (%d)", args, len(args))

		l := len(args)
		switch {
		case l == 0:
			return nil, fmt.Errorf("zero arguments")
		case l%2 == 1:
			return nil, fmt.Errorf("odd number of arguments")
		}

		res := make([]interface{}, l/2)
		for i := 0; i < l; i += 2 {
			res[i/2] = []interface{}{args[i], args[i+1]}
		}
		return res, nil
	}

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()

		script := strings.TrimSpace(`
def check(rows):
	return check_context(rows, {})

def check_context(rows, context):
    return [{"summary": repr(pairs(*rows)), "severity": "notice"}]
		`) + "\n"

		data := []map[string]interface{}{
			{"foo": "bar"},
			{"foo": "baz"},
		}

		addToFuzzCorpus(t.Name(), script, data)
		input, err := PrepareInput(data)
		require.NoError(t, err)

		env, err := NewEnv(t.Name(), script, map[string]GoFunc{"pairs": pairs})
		require.NoError(t, err)

		res, err := env.Run("id", input, nil, t.Log)
		require.NoError(t, err)
		expected := []check.Result{{
			Summary:  `[[{"foo": "bar"}, {"foo": "baz"}]]`,
			Severity: common.Notice,
		}}
		assert.Equal(t, expected, res)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		script := strings.TrimSpace(`
def check(rows):
	return check_context(rows, {})

def check_context(rows, context):
    return [{"summary": repr(pairs(*rows)), "severity": "notice"}]
		`) + "\n"

		data := []map[string]interface{}{}
		addToFuzzCorpus(t.Name(), script, data)
		input, err := PrepareInput(data)
		require.NoError(t, err)

		env, err := NewEnv(t.Name(), script, map[string]GoFunc{"pairs": pairs})
		require.NoError(t, err)

		_, err = env.Run("id", input, nil, t.Log)
		expected := strings.TrimSpace(`
thread id: failed to execute function check_context: pairs: zero arguments
Traceback (most recent call last):
  TestRegisterFunc/Error:5:35: in check_context
  <builtin>: in pairs
		`) + "\n"
		assert.EqualError(t, err, expected)
	})

	t.Run("Kwargs", func(t *testing.T) {
		t.Parallel()

		script := strings.TrimSpace(`
def check(rows):
	return check_context(rows, {})

def check_context(rows, context):
    return [{"summary": repr(pairs(rows=rows)), "severity": "notice"}]
		`) + "\n"

		addToFuzzCorpus(t.Name(), script, nil)
		env, err := NewEnv(t.Name(), script, map[string]GoFunc{"pairs": pairs})
		require.NoError(t, err)

		_, err = env.Run("id", nil, nil, t.Log)
		expected := strings.TrimSpace(`
thread id: failed to execute function check_context: pairs: kwargs are not supported
Traceback (most recent call last):
  TestRegisterFunc/Kwargs:5:35: in check_context
  <builtin>: in pairs
		`) + "\n"
		assert.EqualError(t, err, expected)
	})
}

func TestRegisterAdditionalContext(t *testing.T) {
	t.Parallel()

	concat := func(args ...interface{}) (interface{}, error) {
		l := len(args)
		if l == 0 {
			return nil, fmt.Errorf("zero arguments")
		}

		res := ""
		for i := 0; i < l; i++ {
			row := args[i].(map[string]interface{}) //nolint:forcetypeassert
			for k, v := range row {
				res += fmt.Sprintf("%s:%s", k, v)
			}
		}
		return res, nil
	}

	script := strings.TrimSpace(`
def check_context(rows, context):
	concat = context.get("concat_rows", fail)
	return [{"summary": concat(*rows), "severity": "notice"}]
		`) + "\n"

	data := []map[string]interface{}{
		{"foo": "bar"},
		{"foo": "baz"},
	}

	addToFuzzCorpus(t.Name(), script, data)
	input, err := PrepareInput(data)
	require.NoError(t, err)

	env, err := NewEnv(t.Name(), script, nil)
	require.NoError(t, err)

	contextFuncs := map[string]GoFunc{
		"concat_rows": GoFunc(concat),
	}
	res, err := env.Run("id", input, contextFuncs, t.Log)
	require.NoError(t, err)
	expected := []check.Result{{
		Summary:  `foo:barfoo:baz`,
		Severity: common.Notice,
	}}
	assert.Equal(t, expected, res)
}

func TestCheckGlobals(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		check  *check.Check
		errStr string
	}{
		{
			name: "invalid script",
			check: &check.Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        check.MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "return 1",
			},
			errStr: ":1:1: return statement not within a function",
		},
		{
			name: "missing check function",
			check: &check.Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        check.MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script: strings.TrimSpace(`
def check_context(rows, context):
    """Check Description"""
    pass
                `),
			},
			errStr: "test_check: no `check` function found",
		},
		{
			name: "missing check_context function",
			check: &check.Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        check.MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script: strings.TrimSpace(`
def check(rows):
    pass
                `),
			},
			errStr: "test_check: no `check_context` function found",
		},
		{
			name: "valid script",
			check: &check.Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        check.MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script: strings.TrimSpace(`
def check(rows):
    return check_context(rows, {})

def check_context(rows, context):
    """Check Description"""
    pass
                `),
			},
			errStr: "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := CheckGlobals(tt.check, nil)
			if tt.errStr != "" {
				assert.EqualError(t, err, tt.errStr)
				return
			}

			assert.NoError(t, err)
		})
	}
}
