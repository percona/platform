package check

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseChecks(t *testing.T) {
	t.Parallel()
	monoDocument := strings.TrimSpace(`
---
checks:
  - version: 1
    name: mysql_check
    summary: MYSQL Check
    description: Description of check.
    advisor: test_advisor
    type: MYSQL_SHOW
    query: VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');
    script: |
        def function1(args):
            pass

  - version: 1
    name: postgresql_check
    summary: MYSQL Check
    description: Description of check.
    advisor: test_advisor
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
    advisor: test_advisor
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
    advisor: test_advisor
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

			cs, err := ParseChecks(bytes.NewReader([]byte(document)), params)
			require.NoError(t, err)

			assert.Len(t, cs, 2)

			assert.Equal(t, "mysql_check", cs[0].Name)
			assert.Equal(t, uint32(1), cs[0].Version)
			assert.Equal(t, MySQLShow, cs[0].Type)
			assert.Equal(t, "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');", cs[0].Query)
			assert.Equal(t, cs[0].Script, "def function1(args):\n    pass\n")

			assert.Equal(t, "postgresql_check", cs[1].Name)
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
    advisor: test_advisor
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
		cs, err := ParseChecks(bytes.NewReader([]byte(data)), params)
		require.NoError(t, err)

		assert.Len(t, cs, 1)

		assert.Equal(t, "mysql_check", cs[0].Name)
		assert.Equal(t, uint32(1), cs[0].Version)
		assert.Equal(t, MySQLShow, cs[0].Type)
		assert.Equal(t, "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');", cs[0].Query)
		assert.Equal(t, cs[0].Script, "def function1(args):\n    pass\n")
	})
}

func TestCheck_CheckValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		check  *Check
		errStr string
	}{
		{
			name: "mysql show",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "",
		}, {
			name: "mysql select",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Type:        MySQLSelect,
				Query:       "id, name FROM table WHERE id=123;",
				Script:      "def func(args): pass",
			},
			errStr: "",
		}, {
			name: "postgresql show",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Type:        PostgreSQLShow,
				Query:       "",
				Script:      "def func(args): pass",
			},
			errStr: "",
		}, {
			name: "postgresql select",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Type:        PostgreSQLSelect,
				Query:       "id, name FROM table WHERE id=123;",
				Script:      "def func(args): pass",
			},
			errStr: "",
		}, {
			name: "mongodb getParameter",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Type:        MongoDBGetParameter,
				Script:      "def func(args): pass",
			},
			errStr: "",
		}, {
			name: "mongodb buildInfo",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Type:        MongoDBBuildInfo,
				Script:      "def func(args): pass",
			},
			errStr: "",
		}, {
			name: "mongodb getCmdLineOpts",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Type:        MongoDBGetCmdLineOpts,
				Script:      "def func(args): pass",
			},
			errStr: "",
		}, {
			name: "mongodb replSetGetStatus",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Type:        MongoDBReplSetGetStatus,
				Script:      "def func(args): pass",
			},
			errStr: "",
		}, {
			name: "mongodb getDiagnosticData",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Type:        MongoDBGetDiagnosticData,
				Script:      "def func(args): pass",
			},
			errStr: "",
		}, {
			name: "clickhouse show",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Type:        "CLICKHOUSE_SHOW",
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "unknown check type: CLICKHOUSE_SHOW",
		}, {
			name: "empty version",
			check: &Check{
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "unexpected version 0",
		}, {
			name: "unknown version",
			check: &Check{
				Version:     123,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "unexpected version 123",
		}, {
			name: "empty name",
			check: &Check{
				Version:     1,
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "invalid check name",
		}, {
			name: "invalid_name",
			check: &Check{
				Version:     1,
				Name:        "test check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "invalid check name",
		}, {
			name: "normal interval",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Interval:    Standard,
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "",
		}, {
			name: "empty interval",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Interval:    "",
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "",
		}, {
			name: "unknown interval",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Interval:    Interval("unknown"),
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "unknown check interval: unknown",
		}, {
			name: "empty tier",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "",
		}, {
			name: "empty type",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Type:        "",
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "check type is empty",
		}, {
			name: "empty query",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Type:        MySQLShow,
				Query:       "",
				Script:      "def func(args): pass",
			},
			errStr: "query is empty",
		}, {
			name: "non empty query for postgresql show",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Type:        PostgreSQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "query should be empty for 'POSTGRESQL_SHOW' type",
		}, {
			name: "non empty query for mongodb get parameter",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Type:        MongoDBGetParameter,
				Query:       "some query",
				Script:      "def func(args): pass",
			},
			errStr: "query should be empty for 'MONGODB_GETPARAMETER' type",
		}, {
			name: "non empty query for mongodb build info",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Type:        MongoDBBuildInfo,
				Query:       "some query",
				Script:      "def func(args): pass",
			},
			errStr: "query should be empty for 'MONGODB_BUILDINFO' type",
		}, {
			name: "non empty query for mongodb get cmd line opts",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Type:        MongoDBGetCmdLineOpts,
				Advisor:     "test_advisor",
				Query:       "some query",
				Script:      "def func(args): pass",
			},
			errStr: "query should be empty for 'MONGODB_GETCMDLINEOPTS' type",
		}, {
			name: "non-empty query for mongodb replSetGetStatus",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Type:        MongoDBReplSetGetStatus,
				Advisor:     "test_advisor",
				Query:       "some query",
				Script:      "def func(args): pass",
			},
			errStr: "query should be empty for 'MONGODB_REPLSETGETSTATUS' type",
		}, {
			name: "non-empty query for mongodb getDiagnosticData",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Type:        MongoDBGetDiagnosticData,
				Query:       "some query",
				Script:      "def func(args): pass",
			},
			errStr: "query should be empty for 'MONGODB_GETDIAGNOSTICDATA' type",
		}, {
			name: "empty script",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "",
			},
			errStr: "check script is empty",
		}, {
			name: "empty summary",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "summary is empty",
		}, {
			name: "empty description",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "",
				Advisor:     "test_advisor",
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "description is empty",
		}, {
			name: "script with tabs",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args):\tpass",
			},
			errStr: "script should use spaces for indentation, not tabs",
		}, {
			name: "missing advisor",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "",
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "advisor name is missing",
		}, {
			name: "mysql family check v2",
			check: &Check{
				Version:     2,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Family:      MySQL,
				Queries: []Query{
					{
						Type:  MySQLShow,
						Query: "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
					},
					{
						Type:  MySQLSelect,
						Query: "id, name FROM table WHERE id=123;",
					},
				},
				Script: "def func(args): pass",
			},
			errStr: "",
		}, {
			name: "postgresql family check v2",
			check: &Check{
				Version:     2,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Family:      PostgreSQL,
				Queries: []Query{
					{
						Type: PostgreSQLShow,
					},
					{
						Type:  PostgreSQLSelect,
						Query: "id, name FROM table WHERE id=123;",
					},
				},
				Script: "def func(args): pass",
			},
			errStr: "",
		}, {
			name: "mongodb family check v2",
			check: &Check{
				Version:     2,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Family:      MongoDB,
				Queries: []Query{
					{
						Type: MongoDBGetCmdLineOpts,
					},
					{
						Type: MongoDBGetParameter,
					},
					{
						Type: MongoDBBuildInfo,
					},
					{
						Type: MongoDBGetDiagnosticData,
					},
					{
						Type: MongoDBReplSetGetStatus,
					},
				},
				Script: "def func(args): pass",
			},
			errStr: "",
		}, {
			name: "metrics instant",
			check: &Check{
				Version:     2,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Family:      MongoDB,
				Queries: []Query{
					{
						Type:  MetricsInstant,
						Query: "instant query",
					},
				},
				Script: "def func(args): pass",
			},
			errStr: "",
		}, {
			name: "metrics instant with params",
			check: &Check{
				Version:     2,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Family:      MongoDB,
				Queries: []Query{
					{
						Type:  MetricsInstant,
						Query: "instant query",
						Parameters: map[Parameter]string{
							Lookback: "5m",
						},
					},
				},
				Script: "def func(args): pass",
			},
			errStr: "",
		}, {
			name: "metrics range",
			check: &Check{
				Version:     2,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Family:      MongoDB,
				Queries: []Query{
					{
						Type:  MetricsRange,
						Query: "range query",
						Parameters: map[Parameter]string{
							Range: "20m",
							Step:  "5m",
						},
					},
				},
				Script: "def func(args): pass",
			},
			errStr: "",
		}, {
			name: "metrics range with all params",
			check: &Check{
				Version:     2,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Family:      MongoDB,
				Queries: []Query{
					{
						Type:  MetricsRange,
						Query: "range query",
						Parameters: map[Parameter]string{
							Lookback: "1h",
							Range:    "20m",
							Step:     "5m",
						},
					},
				},
				Script: "def func(args): pass",
			},
			errStr: "",
		}, {
			name: "metrics range with unknown parameter",
			check: &Check{
				Version:     2,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Family:      MongoDB,
				Queries: []Query{
					{
						Type:  MetricsRange,
						Query: "instant query",
						Parameters: map[Parameter]string{
							Parameter("unknown"): "20m",
							Step:                 "5m",
							Range:                "20m",
						},
					},
				},
				Script: "def func(args): pass",
			},
			errStr: "unsupported parameter 'unknown' for range metris query",
		}, {
			name: "metrics range without required parameter",
			check: &Check{
				Version:     2,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Family:      MongoDB,
				Queries: []Query{
					{
						Type:  MetricsRange,
						Query: "range query",
						Parameters: map[Parameter]string{
							Step: "5m",
						},
					},
				},
				Script: "def func(args): pass",
			},
			errStr: "query parameter 'range' is required for metrics range queries",
		}, {
			name: "metrics instant with bad lookback parameter value",
			check: &Check{
				Version:     2,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Family:      MongoDB,
				Queries: []Query{
					{
						Type:  MetricsInstant,
						Query: "instant query",
						Parameters: map[Parameter]string{
							Lookback: "invalid",
						},
					},
				},
				Script: "def func(args): pass",
			},
			errStr: "failed to parse loopback parameter value 'invalid', it should be a duration: time: invalid duration \"invalid\"",
		}, {
			name: "metrics instant with empty query",
			check: &Check{
				Version:     2,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Family:      MongoDB,
				Queries: []Query{
					{
						Type:  MetricsInstant,
						Query: "",
					},
				},
				Script: "def func(args): pass",
			},
			errStr: "query is empty",
		}, {
			name: "metrics range with empty query",
			check: &Check{
				Version:     2,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Family:      MongoDB,
				Queries: []Query{
					{
						Type:  MetricsRange,
						Query: "",
						Parameters: map[Parameter]string{
							Range: "20m",
							Step:  "5m",
						},
					},
				},
				Script: "def func(args): pass",
			},
			errStr: "query is empty",
		}, {
			name: "unsupported query type for given family",
			check: &Check{
				Version:     2,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Family:      MySQL,
				Queries: []Query{
					{
						Type:  MySQLShow,
						Query: "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
					},
					{
						Type:  PostgreSQLSelect,
						Query: "id, name FROM table WHERE id=123;",
					},
				},
				Script: "def func(args): pass",
			},
			errStr: "unsupported query type 'POSTGRESQL_SELECT' for mySQL family",
		}, {
			name: "clickhouse query type is valid for MySQL",
			check: &Check{
				Version:     2,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Category:    "test",
				Family:      MySQL,
				Queries: []Query{
					{
						Type:  MySQLShow,
						Query: "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
					},
					{
						Type:  ClickHouseSelect,
						Query: "id, name FROM metrics WHERE id=123;",
					},
				},
				Script: "def func(args): pass",
			},
			errStr: "",
		}, {
			name: "clickhouse query type is valid for MongoDB",
			check: &Check{
				Version:     2,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Category:    "test",
				Family:      MongoDB,
				Queries: []Query{
					{
						Type:  MongoDBGetParameter,
						Query: "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
					},
					{
						Type:  ClickHouseSelect,
						Query: "id, name FROM metrics WHERE id=123;",
					},
				},
				Script: "def func(args): pass",
			},
			errStr: "",
		}, {
			name: "clickhouse query type is valid for PostgreSQL",
			check: &Check{
				Version:     2,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Tiers:       []common.Tier{common.Anonymous},
				Category:    "test",
				Family:      PostgreSQL,
				Queries: []Query{
					{
						Type:  ClickHouseSelect,
						Query: "id, name FROM metrics WHERE id=123;",
					},
				},
				Script: "def func(args): pass",
			},
			errStr: "",
		}, {
			name: "missing queries",
			check: &Check{
				Version:     2,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Family:      MySQL,
				Queries:     []Query{},
				Script:      "def func(args): pass",
			},
			errStr: "check should have at least one query",
		}, {
			name: "unknown check family",
			check: &Check{
				Version:     2,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Family:      Family("unknown"),
				Queries: []Query{
					{
						Type:  MySQLShow,
						Query: "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
					},
				},
				Script: "def func(args): pass",
			},
			errStr: "unknown check family: unknown",
		}, {
			name: "mixing check format v1 with family field from v2",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Type:        MySQLShow,
				Family:      MySQL,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Script:      "def func(args): pass",
			},
			errStr: "field 'family' is part of check format version 2 and can't be used in version 1",
		}, {
			name: "mixing check format v1 with queries field from v2",
			check: &Check{
				Version:     1,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Type:        MySQLShow,
				Query:       "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
				Queries: []Query{
					{
						Type:  MySQLShow,
						Query: "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
					},
				}, Script: "def func(args): pass",
			},
			errStr: "field 'queries' is part of check format version 2 and can't be used in version 1",
		}, {
			name: "mixing check format v1 with type field from v2",
			check: &Check{
				Version:     2,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Family:      MySQL,
				Type:        MySQLShow,
				Queries: []Query{
					{
						Type:  MySQLShow,
						Query: "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
					},
				},
				Script: "def func(args): pass",
			},
			errStr: "field 'type' is part of check format version 1 and can't be used in version 2",
		}, {
			name: "mixing check format v2 with query field from v1",
			check: &Check{
				Version:     2,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "test_advisor",
				Family:      MySQL,
				Query:       "some query",
				Queries: []Query{
					{
						Type:  MySQLShow,
						Query: "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
					},
				},
				Script: "def func(args): pass",
			},
			errStr: "field 'query' is part of check format version 1 and can't be used in version 2",
		}, {
			name: "missing advisor check v2",
			check: &Check{
				Version:     2,
				Name:        "test_check",
				Summary:     "Test Check",
				Description: "Check Description",
				Advisor:     "",
				Family:      MySQL,
				Queries: []Query{
					{
						Type:  MySQLShow,
						Query: "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');",
					},
					{
						Type:  MySQLSelect,
						Query: "id, name FROM table WHERE id=123;",
					},
				},
				Script: "def func(args): pass",
			},
			errStr: "advisor name is missing",
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

const (
	data      = `random data`
	publicKey = `RWS3wNj+cjvpdKiPgiiqRsbOEPtTP++3Me64W3txOXOtoeplPQciXOu/`
	signature = `untrusted comment: signature from minisign secret key
RWS3wNj+cjvpdJ6ZzxAlsmfz6WGJHICa8umTeLyqfA/ZYKPeJWmhDP+Sn2qf3kgotbQ05eqv4ezvkPiq+QK65ZumPm/Zpk0BtAQ=
trusted comment: timestamp:1638271463	file:data
Ev7cLRh4ftaZMS+97g3U3/9Ic4QpNGtB55AFa33Bwf0V6psv69U7K3nzq+2/j2tz8EqqXCE0iAlAnUxmU9EzDw==
`
	signatureHashed = `untrusted comment: signature from minisign secret key
RUS3wNj+cjvpdG9sn3QKgnnJW2ZUdcOYI+7czEllp3x6ZBJwgbxZS94t8bNYRA5++4p67+JpIm6bn9eMO7b2BbJRUGZVggJxgg8=
trusted comment: timestamp:1638281678	file:data	hashed
Q4aSH3jbkkgKaPlFfL4J9SSKVtxT37v8+o1pXrGN4banCESh1o61qiI42x2wVrJpSOz7BOgjkmP2nbaK/oihBQ==
`
)

// Private key that was used to sign test data.
/*
untrusted comment: minisign encrypted secret key
RWRTY0Iyr0t5TaUWsOUUhtYhUm+QKu+jch5Q/KEKoWIZFi7GcFsAAAACAAAAAAAAAEAAAAAAI+0TaT6z3ylgJ1Wgkf2WDDkXe3kC/acK0dW5vm0TV6zRC1Sfzeoqd+WJleSHYZgr6VPV7VOpgypMw/duwW+69ZeCwsUyTXUmW7NUKWPo41M7t0NSDyhKkGKg8FMONV3Ly29Eb9seK8I=
*/

func TestCheck_Verify(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		data      string
		key       string
		signature string
		errStr    string
	}{
		{
			name:      "valid non hashed",
			data:      data,
			key:       publicKey,
			signature: signature,
			errStr:    "",
		},
		{
			name:      "valid hashed",
			data:      data,
			key:       publicKey,
			signature: signatureHashed,
			errStr:    "",
		},
		{
			name:      "invalid key algorithm",
			data:      data,
			key:       "INVALID+cjvpdKiPgiiqRsbOEPtTP++3Me64W3txOXOtoeplPQciXOu/",
			signature: signatureHashed,
			errStr:    "unsupported key algorithm",
		},
		{
			name: "invalid signature algorithm",
			data: data,
			key:  publicKey,
			signature: `untrusted comment: signature from minisign secret key
INVALID+cjvpdJ6ZzxAlsmfz6WGJHICa8umTeLyqfA/ZYKPeJWmhDP+Sn2qf3kgotbQ05eqv4ezvkPiq+QK65ZumPm/Zpk0BtAQ=
trusted comment: timestamp:1638271463	file:data
Ev7cLRh4ftaZMS+97g3U3/9Ic4QpNGtB55AFa33Bwf0V6psv69U7K3nzq+2/j2tz8EqqXCE0iAlAnUxmU9EzDw==`,
			errStr: "unsupported signature algorithm",
		},
		{
			name:      "incompatible key identifiers",
			data:      data,
			key:       "RWS3wNj+cINVdKiPgiiqRsbOEPtTP++3Me64W3txOXOtoeplPQciXOu/",
			signature: signature,
			errStr:    "incompatible key identifiers",
		},
		{
			name: "invalid signature",
			data: data,
			key:  publicKey,
			signature: `untrusted comment: signature from minisign secret key
		RWS3wNj+cjvpdJ6ZzxAlsmfz6WGJHICa8umTeLyqfA/ZYKPeJWmhDP+Sn2qf3kgotbQ05eqv4ezvkPiq+QK+INVALID+INVALID=
		trusted comment: timestamp:1638271463	file:data
		Ev7cLRh4ftaZMS+97g3U3/9Ic4QpNGtB55AFa33Bwf0V6psv69U7K3nzq+2/j2tz8EqqXCE0iAlAnUxmU9EzDw==`,
			errStr: "invalid signature",
		},
		{
			name: "invalid global signature",
			data: data,
			key:  publicKey,
			signature: `untrusted comment: signature from minisign secret key
RWS3wNj+cjvpdJ6ZzxAlsmfz6WGJHICa8umTeLyqfA/ZYKPeJWmhDP+Sn2qf3kgotbQ05eqv4ezvkPiq+QK65ZumPm/Zpk0BtAQ=
trusted comment: timestamp:1638271463	file:data
Ev7cLRh4ftaZMS+97g3U3/9Ic4QpNGtB55AFa33Bwf0V6psv69U7K3nzq+2/j2tz8EqqXC+INVALID+INVALID==`,
			errStr: "invalid global signature",
		},
		{
			name: "invalid trusted comment",
			data: data,
			key:  publicKey,
			signature: `untrusted comment: signature from minisign secret key
RWS3wNj+cjvpdJ6ZzxAlsmfz6WGJHICa8umTeLyqfA/ZYKPeJWmhDP+Sn2qf3kgotbQ05eqv4ezvkPiq+QK65ZumPm/Zpk0BtAQ=
trusted comment: timestamp:1638271463	file:INVALID COMMENT
Ev7cLRh4ftaZMS+97g3U3/9Ic4QpNGtB55AFa33Bwf0V6psv69U7K3nzq+2/j2tz8EqqXCE0iAlAnUxmU9EzDw==`,
			errStr: "invalid global signature",
		},
		{
			name:      "empty data",
			data:      "",
			key:       publicKey,
			signature: signature,
			errStr:    "invalid signature",
		},
		{
			name:      "empty signature",
			data:      data,
			key:       publicKey,
			signature: "",
			errStr:    "incomplete signature",
		},
		{
			name:      "empty key",
			data:      data,
			key:       "",
			signature: signature,
			errStr:    "invalid public key",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := Verify([]byte(tt.data), tt.key, tt.signature)

			if tt.errStr != "" {
				assert.EqualError(t, err, tt.errStr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
