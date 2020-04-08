package checked

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheck_ParseChecks(t *testing.T) {
	t.Run("singleFile", func(t *testing.T) {
		var data = `---
checks:
  - type: MYSQL_SHOW
    query: VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');
    script: |
        def function1(args):
            pass
`

		cs, err := ParseChecks(bytes.NewReader([]byte(data)))
		require.NoError(t, err)

		assert.Len(t, cs, 1)
		assert.Equal(t, "MYSQL_SHOW", cs[0].Type)
		assert.Equal(t, "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');", cs[0].Query)
		assert.Equal(t, cs[0].Script, "def function1(args):\n    pass\n")
	})

	t.Run("multipleFiles", func(t *testing.T) {
		var data = `---
checks:
  - type: MYSQL_SHOW
    query: VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');
    script: |
        def function1(args):
            pass
---
checks:
  - type: POSTGRESQL_SELECT
    query: id, name FROM table WHERE id=123;
    script: |
        def function2(args):
            pass
`
		cs, err := ParseChecks(bytes.NewReader([]byte(data)))
		require.NoError(t, err)

		assert.Len(t, cs, 2)

		assert.Equal(t, "MYSQL_SHOW", cs[0].Type)
		assert.Equal(t, "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');", cs[0].Query)
		assert.Equal(t, cs[0].Script, "def function1(args):\n    pass\n")

		assert.Equal(t, "POSTGRESQL_SELECT", cs[1].Type)
		assert.Equal(t, "id, name FROM table WHERE id=123;", cs[1].Query)
		assert.Equal(t, cs[1].Script, "def function2(args):\n    pass\n")
	})
}

func TestCheck_ParseCheck(t *testing.T) {
	data := `type: MYSQL_SHOW
query: VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');
script: |
    def function1(args):
        pass
`

	c, err := ParseCheck([]byte(data))
	require.NoError(t, err)

	assert.Equal(t, "MYSQL_SHOW", c.Type)
	assert.Equal(t, "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');", c.Query)
	assert.Equal(t, c.Script, "def function1(args):\n    pass\n")
}

func TestCheck_CheckValidate(t *testing.T) {
	tests := []struct {
		name   string
		check  *Check
		errStr string
	}{
		{
			name:   "mysql_show",
			check:  &Check{Type: "MYSQL_SHOW", Query: "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');", Script: "def func(args): pass"},
			errStr: "",
		},
		{
			name:   "mysql_select",
			check:  &Check{Type: "MYSQL_SELECT", Query: "id, name FROM table WHERE id=123;", Script: "def func(args): pass"},
			errStr: "",
		},
		{
			name:   "postgresql_show",
			check:  &Check{Type: "POSTGRESQL_SHOW", Query: "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');", Script: "def func(args): pass"},
			errStr: "",
		},
		{
			name:   "postgresql_select",
			check:  &Check{Type: "POSTGRESQL_SELECT", Query: "id, name FROM table WHERE id=123;", Script: "def func(args): pass"},
			errStr: "",
		},
		{
			name:   "mongodb_get_parameter",
			check:  &Check{Type: "MONGODB_GETPARAMETER", Query: "\"saslHostName\" : 1", Script: "def func(args): pass"},
			errStr: "",
		},
		{
			name:   "clickhouse_show",
			check:  &Check{Type: "CLICKHOUSE_SHOW", Query: "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');", Script: "def func(args): pass"},
			errStr: "unknown check type: CLICKHOUSE_SHOW",
		},
		{
			name:   "empty_type",
			check:  &Check{Type: "", Query: "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');", Script: "def func(args): pass"},
			errStr: "check type is empty",
		},
		{
			name:   "empty_query",
			check:  &Check{Type: "MYSQL_SHOW", Query: "", Script: "def func(args): pass"},
			errStr: "check query is empty",
		},
		{
			name:   "empty_script",
			check:  &Check{Type: "MYSQL_SHOW", Query: "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');", Script: ""},
			errStr: "check script is empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.check.Validate()

			if tt.errStr != "" {
				assert.EqualError(t, err, tt.errStr)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestCheck_ParseResults(t *testing.T) {
	t.Run("singleFile", func(t *testing.T) {
		var data = `---
results:
  - status: FAIL
    message: something bad happened
`
		cs, err := ParseResults(bytes.NewReader([]byte(data)))
		require.NoError(t, err)

		assert.Len(t, cs, 1)
		assert.Equal(t, "FAIL", cs[0].Status)
		assert.Equal(t, "something bad happened", cs[0].Message)
	})

	t.Run("multipleFiles", func(t *testing.T) {
		var data = `---
results:
  - status: FAIL
    message: something bad happened
---
results:
  - status: SUCCESS
    message: 
`
		cs, err := ParseResults(bytes.NewReader([]byte(data)))
		require.NoError(t, err)

		assert.Len(t, cs, 2)

		assert.Equal(t, "FAIL", cs[0].Status)
		assert.Equal(t, "something bad happened", cs[0].Message)

		assert.Equal(t, "SUCCESS", cs[1].Status)
		assert.Empty(t,  cs[1].Message)
	})
}

func TestCheck_ParseResult(t *testing.T) {
	data := `status: FAIL
message: something bad happened
`

	r, err := ParseResult([]byte(data))
	require.NoError(t, err)

	assert.Equal(t, "FAIL", r.Status)
	assert.Equal(t, "something bad happened", r.Message)
}

func TestCheck_ResultValidate(t *testing.T) {
	tests := []struct {
		name   string
		result *Result
		errStr string
	}{
		{
			name:   "success_result_without_message",
			result: &Result{Status: Success, Message: ""},
			errStr: "",
		},
		{
			name:   "success_result_with_message",
			result: &Result{Status: Success, Message: "everything is fine!"},
			errStr: "",
		},
		{
			name:   "failed_result_with_message",
			result: &Result{Status: Fail, Message: "something bad happened!"},
			errStr: "",
		},
		{
			name:   "failed_result_without_message",
			result: &Result{Status: Fail, Message: ""},
			errStr: "failed check result should have message",
		},
		{
			name:   "empty_status",
			result: &Result{Status: "", Message: ""},
			errStr: "result status is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.result.Validate()

			if tt.errStr != "" {
				assert.EqualError(t, err, tt.errStr)
				return
			}

			assert.NoError(t, err)
		})
	}
}
