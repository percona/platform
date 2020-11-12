package check

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.starlark.net/starlark"

	"github.com/percona-platform/platform/pkg/common"
)

func TestCheck_Parse(t *testing.T) {
	t.Parallel()
	monoDocument := strings.TrimSpace(`
---
checks:
  - version: 1
    name: mysql_check
    summary: MYSQL Check
    description: Description of check.
    tiers: [anonymous]
    type: MYSQL_SHOW
    query: VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');
    script: |
        def function1(args):
            pass

  - version: 1
    name: postgresql_check
    summary: MYSQL Check
    description: Description of check.
    tiers: [anonymous]
    type: POSTGRESQL_SELECT
    query: id, name FROM table WHERE id=123;
    script: |
        def function2(args):
            pass
`)

	multiDocument := strings.TrimSpace(`
---
checks:
  - version: 1
    name: mysql_check
    summary: MYSQL Check
    description: Description of check.
    tiers: [anonymous]
    type: MYSQL_SHOW
    query: VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');
    script: |
        def function1(args):
            pass
---
checks:
  - version: 1
    name: postgresql_check
    summary: PostgreSQL Check
    description: Description of check.
    tiers: [anonymous]
    type: POSTGRESQL_SELECT
    query: id, name FROM table WHERE id=123;
    script: |
        def function2(args):
            pass
`)

	params := &ParseParams{
		DisallowUnknownFields: true,
		DisallowInvalidChecks: true,
	}

	for name, document := range map[string]string{"mono-document": monoDocument, "multi-document": multiDocument} {
		name, document := name, document
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			cs, err := Parse(bytes.NewReader([]byte(document)), params)
			require.NoError(t, err)

			assert.Len(t, cs, 2)

			assert.Equal(t, "mysql_check", cs[0].Name)
			assert.Equal(t, []common.Tier{common.Anonymous}, cs[0].Tiers)
			assert.Equal(t, uint32(1), cs[0].Version)
			assert.Equal(t, MySQLShow, cs[0].Type)
			assert.Equal(t, "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');", cs[0].Query)
			assert.Equal(t, cs[0].Script, "def function1(args):\n    pass\n")

			assert.Equal(t, "postgresql_check", cs[1].Name)
			assert.Equal(t, []common.Tier{common.Anonymous}, cs[0].Tiers)
			assert.Equal(t, uint32(1), cs[1].Version)
			assert.Equal(t, PostgreSQLSelect, cs[1].Type)
			assert.Equal(t, "id, name FROM table WHERE id=123;", cs[1].Query)
			assert.Equal(t, cs[1].Script, "def function2(args):\n    pass")
		})
	}

	t.Run("skipInvalid", func(t *testing.T) {
		t.Parallel()
		data := strings.TrimSpace(`
---
checks:
  - version: 1
    name: mysql_check
    summary: MYSQL Check
    description: Description of check.
    tiers: [anonymous]
    type: MYSQL_SHOW
    query: VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');
    script: |
        def function1(args):
            pass

  - version: 2
`)

		params := &ParseParams{
			DisallowUnknownFields: true,
			DisallowInvalidChecks: false,
		}
		cs, err := Parse(bytes.NewReader([]byte(data)), params)
		require.NoError(t, err)

		assert.Len(t, cs, 1)

		assert.Equal(t, "mysql_check", cs[0].Name)
		assert.Equal(t, []common.Tier{common.Anonymous}, cs[0].Tiers)
		assert.Equal(t, uint32(1), cs[0].Version)
		assert.Equal(t, MySQLShow, cs[0].Type)
		assert.Equal(t, "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');", cs[0].Query)
		assert.Equal(t, cs[0].Script, "def function1(args):\n    pass\n")
	})

	t.Run("missing tiers", func(t *testing.T) {
		t.Parallel()
		data := strings.TrimSpace(`
---
checks:
  - version: 1
    name: mysql_check
    summary: MYSQL Check
    description: Description of check.
    type: MYSQL_SHOW
    query: VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');
    script: |
        def function1(args):
            pass
`)

		params := &ParseParams{
			DisallowUnknownFields: true,
			DisallowInvalidChecks: true,
		}
		cs, err := Parse(bytes.NewReader([]byte(data)), params)
		require.NoError(t, err)

		assert.Len(t, cs, 1)
		assert.Nil(t, cs[0].Tiers)
	})

	t.Run("null tiers", func(t *testing.T) {
		t.Parallel()
		data := strings.TrimSpace(`
---
checks:
  - version: 1
    name: mysql_check
    summary: MYSQL Check
    description: Description of check.
    tiers: null
    type: MYSQL_SHOW
    query: VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');
    script: |
        def function1(args):
            pass
`)

		params := &ParseParams{
			DisallowUnknownFields: true,
			DisallowInvalidChecks: true,
		}
		cs, err := Parse(bytes.NewReader([]byte(data)), params)
		require.NoError(t, err)

		assert.Len(t, cs, 1)

		assert.Equal(t, "mysql_check", cs[0].Name)
		assert.Len(t, cs[0].Tiers, 0)
		assert.Equal(t, uint32(1), cs[0].Version)
		assert.Equal(t, MySQLShow, cs[0].Type)
		assert.Equal(t, "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');", cs[0].Query)
		assert.Equal(t, cs[0].Script, "def function1(args):\n    pass")
	})

	t.Run("zero tiers", func(t *testing.T) {
		t.Parallel()
		data := strings.TrimSpace(`
---
checks:
  - version: 1
    name: mysql_check
    summary: MYSQL Check
    description: Description of check.
    tiers: []
    type: MYSQL_SHOW
    query: VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');
    script: |
        def function1(args):
            pass
`)

		params := &ParseParams{
			DisallowUnknownFields: true,
			DisallowInvalidChecks: true,
		}
		cs, err := Parse(bytes.NewReader([]byte(data)), params)
		require.NoError(t, err)

		assert.Len(t, cs, 1)

		assert.Equal(t, "mysql_check", cs[0].Name)
		assert.Len(t, cs[0].Tiers, 0)
		assert.Equal(t, uint32(1), cs[0].Version)
		assert.Equal(t, MySQLShow, cs[0].Type)
		assert.Equal(t, "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');", cs[0].Query)
		assert.Equal(t, cs[0].Script, "def function1(args):\n    pass")
	})

	t.Run("duplicate tiers", func(t *testing.T) {
		t.Parallel()
		data := strings.TrimSpace(`
---
checks:
  - version: 1
    name: mysql_check
    summary: MYSQL Check
    description: Description of check.
    tiers: [anonymous, anonymous]
    type: MYSQL_SHOW
    query: VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');
    script: |
        def function1(args):
            pass
`)

		params := &ParseParams{
			DisallowUnknownFields: true,
			DisallowInvalidChecks: true,
		}
		_, err := Parse(bytes.NewReader([]byte(data)), params)
		require.EqualError(t, err, "duplicate tier: \"anonymous\"")
	})
}

func TestCheck_GetDocstring(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		check     *Check
		docstring string
		errStr    string
	}{
		{
			name: "invalid script",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "return 1",
			},
			docstring: "",
			errStr:    ":1:1: return statement not within a function",
		},
		{
			name: "missing check function",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script: strings.TrimSpace(`
def check_context(rows, context):
    """Check Description"""
    pass
                `),
			},
			docstring: "",
			errStr:    "test_check: no `check` function found",
		},
		{
			name: "missing check_context function",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script: strings.TrimSpace(`
def check(rows):
    pass
                `),
			},
			docstring: "",
			errStr:    "test_check: no `check_context` function found",
		},
		{
			name: "missing docstring",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script: strings.TrimSpace(`
def check(rows):
    return check_context(rows, {})
           
def check_context(rows, context):
    pass
                `),
			},
			docstring: "",
			errStr:    "test_check: `check_context` function should have docstring",
		},
		{
			name: "valid script",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script: strings.TrimSpace(`
def check(rows):
    return check_context(rows, {})

def check_context(rows, context):
    """Check Description"""
    pass
                `),
			},
			docstring: "Check Description",
			errStr:    "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			doc, err := tt.check.GetDocstring(starlark.StringDict{})
			assert.Equal(t, tt.docstring, doc)

			if tt.errStr != "" {
				assert.EqualError(t, err, tt.errStr)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestCheck_CheckValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		check  *Check
		errStr string
	}{
		{
			name: "mysql_show",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "",
		},
		{
			name: "mysql_select",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        MySQLSelect,
				Query:       "id, name FROM table WHERE id=123;",
				Script:      "def func(args): pass",
			},
			errStr: "",
		},
		{
			name: "postgresql_show",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        PostgreSQLShow,
				Query:       "",
				Script:      "def func(args): pass",
			},
			errStr: "",
		},
		{
			name: "postgresql_select",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        PostgreSQLSelect,
				Query:       "id, name FROM table WHERE id=123;",
				Script:      "def func(args): pass",
			},
			errStr: "",
		},
		{
			name: "mongodb_get_parameter",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        MongoDBGetParameter,
				Script:      "def func(args): pass",
			},
			errStr: "",
		},
		{
			name: "mongodb_build_info",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        MongoDBBuildInfo,
				Script:      "def func(args): pass",
			},
			errStr: "",
		},
		{
			name: "mongodb_get_cmd_line_opts",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        MongoDBGetCmdLineOpts,
				Script:      "def func(args): pass",
			},
			errStr: "",
		},
		{
			name: "clickhouse_show",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        "CLICKHOUSE_SHOW",
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "unknown check type: CLICKHOUSE_SHOW",
		},
		{
			name: "empty_version",
			check: &Check{
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "unexpected version 0",
		},
		{
			name: "empty_name",
			check: &Check{
				Version:     1,
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "invalid check name",
		},
		{
			name: "invalid_name",
			check: &Check{
				Version:     1,
				Name:        "test check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "invalid check name",
		},
		{
			name: "empty_tier",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "",
		},
		{
			name: "invalid_tier",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{"invalid"},
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "unknown check tier: \"invalid\"",
		},
		{
			name: "empty_type",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        "",
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "check type is empty",
		},
		{
			name: "empty_query",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        MySQLShow,
				Query:       "",
				Script:      "def func(args): pass",
			},
			errStr: "check query is empty",
		},
		{
			name: "non_empty_query_for_postgresql_show",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        PostgreSQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "POSTGRESQL_SHOW check type should have empty query",
		},
		{
			name: "non_empty_query_for_mongodb_get_parameter",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        MongoDBGetParameter,
				Query:       "some query",
				Script:      "def func(args): pass",
			},
			errStr: "MONGODB_GETPARAMETER check type should have empty query",
		},
		{
			name: "non_empty_query_for_mongodb_build_info",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        MongoDBBuildInfo,
				Query:       "some query",
				Script:      "def func(args): pass",
			},
			errStr: "MONGODB_BUILDINFO check type should have empty query",
		},
		{
			name: "non_empty_query_for_mongodb_get_cmd_line_opts",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        MongoDBGetCmdLineOpts,
				Query:       "some query",
				Script:      "def func(args): pass",
			},
			errStr: "MONGODB_GETCMDLINEOPTS check type should have empty query",
		},
		{
			name: "empty_script",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "",
			},
			errStr: "check script is empty",
		},
		{
			name: "empty_summary",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "summary is empty",
		},
		{
			name: "empty_summary",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "description is empty",
		},
		{
			name: "script_with_tabs",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args):\tpass",
			},
			errStr: "script should use spaces for indentation, not tabs",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.check.Validate()

			if tt.errStr != "" {
				assert.EqualError(t, err, tt.errStr)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestCheck_ResultValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		result *Result
		errStr string
	}{
		{
			name:   "normal",
			result: &Result{Severity: common.Notice, Summary: "some text"},
			errStr: "",
		},
		{
			name:   "unknown_severity",
			result: &Result{Severity: common.Severity(123), Summary: "some text"},
			errStr: "unknown severity level: Severity(123)",
		},
		{
			name:   "unhandled_severity",
			result: &Result{Severity: common.Info, Summary: "some text"},
			errStr: "unhandled result severity: info",
		},
		{
			name:   "empty_summary",
			result: &Result{Severity: common.Notice},
			errStr: "summary is empty",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.result.Validate()

			if tt.errStr != "" {
				assert.EqualError(t, err, tt.errStr)
				return
			}

			assert.NoError(t, err)
		})
	}
}

const data = `random data`

const publicKey = `RWRQmBOLeYzAeuR2L6L1GJN9qTR8ceQrawtijPTQkVbf3LJsrLeUjQcL`

const signature = `untrusted comment: signature from minisign secret key
RWRQmBOLeYzAetS6fGVWAvzwCgDuo/zNlvdOrClAvjCUSMLnUimp6NQd1L+x77HZa0kEB7ei+K9lW+W4hIf1D8gRNm+cdQr7dgk=
trusted comment: timestamp:1586854934	file:data
WXAxVyC6G82QuXtGlJZzLWoVmw8QNWks2T6RfXo8F9oKjI+sPbBf0ZOBWD2hXKFBCo5pKPSJiaVeI4G36OlEAw==
`

func TestCheck_Verify(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		err := Verify([]byte(data), publicKey, signature)
		require.NoError(t, err)
	})

	t.Run("invalid signature", func(t *testing.T) {
		err := Verify([]byte(data), publicKey, strings.TrimSpace(`
untrusted comment: signature from minisign secret key
RWRQmBOLeYzAetS6fGVWAvzwCgDuo/zNlvdOrClAvjCUSMLnUimp6NQd1L+f3fHZa0kEB7ei+K9lW+W4hIf+INVALID+INVALID=
trusted comment: timestamp:1586854934	file:data
WXAxVyC6G82QuXtGlJZzLWoVmw8QNWks2T6RfXo8F9oKjI+sPbBf0ZOBWD2hXKFBCo5pKPSJiaVeI4G36OlEAw==`))

		assert.EqualError(t, err, "invalid signature")
	})

	t.Run("invalid global signature", func(t *testing.T) {
		err := Verify([]byte(data), publicKey, strings.TrimSpace(`
untrusted comment: signature from minisign secret key
RWRQmBOLeYzAetS6fGVWAvzwCgDuo/zNlvdOrClAvjCUSMLnUimp6NQd1L+x77HZa0kEB7ei+K9lW+W4hIf1D8gRNm+cdQr7dgk=
trusted comment: timestamp:1586854934	file:data
WXAxVyC6G82QuXtGlJZzLWoVmw8QNWks2veRfXo8F9oKjI+sPbBf0ZOBWD2hXKFBCo5pKP+INVALID+INVALID==`))
		assert.EqualError(t, err, "invalid global signature")
	})

	t.Run("invalid trusted comment", func(t *testing.T) {
		err := Verify([]byte(data), publicKey, strings.TrimSpace(`
untrusted comment: signature from minisign secret key
RWRQmBOLeYzAetS6fGVWAvzwCgDuo/zNlvdOrClAvjCUSMLnUimp6NQd1L+x77HZa0kEB7ei+K9lW+W4hIf1D8gRNm+cdQr7dgk=
trusted comment: timestamp:1586854934	file:INVALID COMMENT
WXAxVyC6G82QuXtGlJZzLWoVmw8QNWks2T6RfXo8F9oKjI+sPbBf0ZOBWD2hXKFBCo5pKPSJiaVeI4G36OlEAw==`))
		assert.EqualError(t, err, "invalid global signature")
	})

	t.Run("invalid public key", func(t *testing.T) {
		err := Verify([]byte("random data"), "RWRQmBOLeYzAeu5FL8f1JMN9qTR8CDfrabdtjPTQ+INVALID+INVALID", signature)
		assert.EqualError(t, err, "invalid signature")
	})

	t.Run("empty data", func(t *testing.T) {
		err := Verify(nil, publicKey, signature)
		assert.EqualError(t, err, "invalid signature")
	})

	t.Run("empty signature", func(t *testing.T) {
		err := Verify([]byte(data), publicKey, "")
		assert.EqualError(t, err, "incomplete signature")
	})

	t.Run("empty key", func(t *testing.T) {
		err := Verify([]byte(data), "", signature)
		assert.EqualError(t, err, "invalid public key")
	})
}
