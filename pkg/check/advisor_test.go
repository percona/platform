package check

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/percona/platform/pkg/common"
)

func TestParseAdvisors(t *testing.T) {
	t.Parallel()

	monoDocument := strings.TrimSpace(`
---
advisors:
  - version: 1
    name: test_advisor
    summary: Test advisors
    description: Test advisor description.
    category: test
    tiers: [anonymous]
    checks:
      - version: 1
        name: mysql_check
        summary: MYSQL Check
        description: Description of advisor.
        advisor: test_advisor
        type: MYSQL_SHOW
        query: VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');
        script: |
            def function1(args):
                pass

  - version: 1
    name: another_test_advisor
    summary: Another test advisors
    description: Another test advisor description.
    category: test
    tiers: [registered]
    checks:
      - version: 1
        name: postgresql_check
        summary: MYSQL Check
        description: Description of advisor.
        advisor: another_test_advisor
        type: POSTGRESQL_SELECT
        query: id, name FROM table WHERE id=123;
        script: |
            def function2(args):
                pass
`)

	multiDocument := strings.TrimSpace(`
---
advisors:
  - version: 1
    name: test_advisor
    summary: Test advisors
    description: Test advisor description.
    category: test
    tiers: [anonymous]
    checks:
      - version: 1
        name: mysql_check
        summary: MYSQL Check
        description: Description of advisor.
        advisor: test_advisor
        type: MYSQL_SHOW
        query: VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');
        script: |
            def function1(args):
                pass
---
advisors:
  - version: 1
    name: another_test_advisor
    summary: Another test advisors
    description: Another test advisor description.
    category: test
    tiers: [registered]
    checks:
      - version: 1
        name: postgresql_check
        summary: MYSQL Check
        description: Description of advisor.
        advisor: another_test_advisor
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

			as, err := ParseAdvisors(bytes.NewReader([]byte(document)), params)
			require.NoError(t, err)

			assert.Len(t, as, 2)

			assert.Equal(t, "test_advisor", as[0].Name)
			assert.Equal(t, "Test advisors", as[0].Summary)
			assert.Equal(t, "Test advisor description.", as[0].Description)
			assert.Equal(t, "test", as[0].Category)
			assert.Equal(t, []common.Tier{common.Anonymous}, as[0].Tiers)
			assert.Len(t, as[0].Checks, 1)

			assert.Equal(t, "mysql_check", as[0].Checks[0].Name)
			assert.Equal(t, uint32(1), as[0].Checks[0].Version)
			assert.Equal(t, MySQLShow, as[0].Checks[0].Type)
			assert.Equal(t, "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');", as[0].Checks[0].Query)
			assert.Equal(t, as[0].Checks[0].Script, "def function1(args):\n    pass\n")

			assert.Equal(t, "another_test_advisor", as[1].Name)
			assert.Equal(t, "Another test advisors", as[1].Summary)
			assert.Equal(t, "Another test advisor description.", as[1].Description)
			assert.Equal(t, "test", as[1].Category)
			assert.Equal(t, []common.Tier{common.Registered}, as[1].Tiers)
			assert.Len(t, as[0].Checks, 1)

			assert.Equal(t, "postgresql_check", as[1].Checks[0].Name)
			assert.Equal(t, uint32(1), as[1].Checks[0].Version)
			assert.Equal(t, PostgreSQLSelect, as[1].Checks[0].Type)
			assert.Equal(t, "id, name FROM table WHERE id=123;", as[1].Checks[0].Query)
			assert.Equal(t, as[1].Checks[0].Script, "def function2(args):\n    pass")
		})
	}

	t.Run("multiple checks", func(t *testing.T) {
		t.Parallel()

		document := strings.TrimSpace(`
---
advisors:
  - version: 1
    name: test_advisor
    summary: Test advisors
    description: Test advisor description.
    category: test
    tiers: [anonymous]
    checks:
      - version: 1
        name: mysql_check
        summary: MYSQL Check
        description: Description of advisor.
        advisor: test_advisor
        type: MYSQL_SHOW
        query: VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');
        script: |
            def function1(args):
                pass

      - version: 1
        name: postgresql_check
        summary: MYSQL Check
        description: Description of advisor.
        advisor: test_advisor
        type: POSTGRESQL_SELECT
        query: id, name FROM table WHERE id=123;
        script: |
            def function2(args):
                pass
`)

		as, err := ParseAdvisors(bytes.NewReader([]byte(document)), params)
		require.NoError(t, err)

		assert.Equal(t, "test_advisor", as[0].Name)
		assert.Equal(t, "Test advisors", as[0].Summary)
		assert.Equal(t, "Test advisor description.", as[0].Description)
		assert.Equal(t, "test", as[0].Category)
		assert.Equal(t, []common.Tier{common.Anonymous}, as[0].Tiers)
		assert.Len(t, as[0].Checks, 2)

		assert.Equal(t, "mysql_check", as[0].Checks[0].Name)
		assert.Equal(t, uint32(1), as[0].Checks[0].Version)
		assert.Equal(t, MySQLShow, as[0].Checks[0].Type)
		assert.Equal(t, "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');", as[0].Checks[0].Query)
		assert.Equal(t, "def function1(args):\n    pass\n", as[0].Checks[0].Script)

		assert.Equal(t, "postgresql_check", as[0].Checks[1].Name)
		assert.Equal(t, uint32(1), as[0].Checks[1].Version)
		assert.Equal(t, PostgreSQLSelect, as[0].Checks[1].Type)
		assert.Equal(t, "id, name FROM table WHERE id=123;", as[0].Checks[1].Query)
		assert.Equal(t, as[0].Checks[1].Script, "def function2(args):\n    pass")
	})

	t.Run("wrong advisor name specified in advisor", func(t *testing.T) {
		t.Parallel()
		document := strings.TrimSpace(`
---
advisors:
  - version: 1
    name: test_advisor
    summary: Test advisors
    description: Test advisor description.
    category: test
    tiers: [anonymous]
    checks:
      - version: 1
        name: mysql_check
        summary: MYSQL Check
        description: Description of advisor.
        advisor: different_advisor
        type: MYSQL_SHOW
        query: VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');
        script: |
            def function1(args):
                pass
`)
		_, err := ParseAdvisors(bytes.NewReader([]byte(document)), params)
		require.EqualError(t, err, "advisor name 'test_advisor' doesn't match name 'different_advisor' specified in corresponding check 'mysql_check'")
	})

	t.Run("normal no tiers specified", func(t *testing.T) {
		t.Parallel()
		data := strings.TrimSpace(`
---
advisors:
  - version: 1
    name: test_advisor
    summary: Test advisors
    description: Test advisor description.
    category: test
`)

		params := &ParseParams{
			DisallowUnknownFields: true,
			DisallowInvalidChecks: true,
		}
		as, err := ParseAdvisors(bytes.NewReader([]byte(data)), params)
		require.NoError(t, err)

		assert.Len(t, as, 1)

		assert.Equal(t, "test_advisor", as[0].Name)
		assert.Equal(t, "Test advisors", as[0].Summary)
		assert.Equal(t, "Test advisor description.", as[0].Description)
		assert.Equal(t, "test", as[0].Category)
		assert.Empty(t, as[0].Tiers)
		assert.Empty(t, as[0].Checks)
	})

	t.Run("null tiers", func(t *testing.T) {
		t.Parallel()
		data := strings.TrimSpace(`
---
advisors:
  - version: 1
    name: test_advisor
    summary: Test advisors
    description: Test advisor description.
    category: test
    tiers: null
`)

		params := &ParseParams{
			DisallowUnknownFields: true,
			DisallowInvalidChecks: true,
		}
		as, err := ParseAdvisors(bytes.NewReader([]byte(data)), params)
		require.NoError(t, err)

		assert.Len(t, as, 1)

		assert.Equal(t, "test_advisor", as[0].Name)
		assert.Equal(t, "Test advisors", as[0].Summary)
		assert.Equal(t, "Test advisor description.", as[0].Description)
		assert.Equal(t, "test", as[0].Category)
		assert.Empty(t, as[0].Tiers)
		assert.Empty(t, as[0].Checks)
	})

	t.Run("zero tiers", func(t *testing.T) {
		t.Parallel()
		data := strings.TrimSpace(`
---
advisors:
  - version: 1
    name: test_advisor
    summary: Test advisors
    description: Test advisor description.
    category: test
    tiers: []
`)

		params := &ParseParams{
			DisallowUnknownFields: true,
			DisallowInvalidChecks: true,
		}
		as, err := ParseAdvisors(bytes.NewReader([]byte(data)), params)
		require.NoError(t, err)

		assert.Len(t, as, 1)

		assert.Equal(t, "test_advisor", as[0].Name)
		assert.Equal(t, "Test advisors", as[0].Summary)
		assert.Equal(t, "Test advisor description.", as[0].Description)
		assert.Equal(t, "test", as[0].Category)
		assert.Empty(t, as[0].Tiers)
		assert.Empty(t, as[0].Checks)
	})

	t.Run("duplicate tiers", func(t *testing.T) {
		t.Parallel()
		data := strings.TrimSpace(`
---
advisors:
  - version: 1
    name: test_advisor
    summary: Test advisors
    description: Test advisor description.
    category: test
    tiers: [anonymous, anonymous]
`)

		params := &ParseParams{
			DisallowUnknownFields: true,
			DisallowInvalidChecks: true,
		}
		_, err := ParseAdvisors(bytes.NewReader([]byte(data)), params)
		require.EqualError(t, err, "duplicate tier: \"anonymous\"")
	})
}

func TestAdvisor_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		advisor *Advisor
		errStr  string
	}{
		{
			name: "normal",
			advisor: &Advisor{
				Version:     1,
				Name:        "test_advisor",
				Summary:     "Test Advisor",
				Description: "Check Description",
				Category:    "test_category",
				Tiers:       []common.Tier{common.Anonymous},
				Checks: []Check{
					{
						Version:     1,
						Name:        "test_check_1",
						Summary:     "Test Check #1",
						Description: "Check Description",
						Advisor:     "test_advisor",
						Type:        MySQLSelect,
						Query:       "id, name FROM table WHERE id=123;",
						Script:      "def func(args): pass",
					},
					{
						Version:     2,
						Name:        "test_check_2",
						Summary:     "Test Check #2",
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
				},
			},
			errStr: "",
		},
		{
			name: "normal",
			advisor: &Advisor{
				Version:     1,
				Name:        "test_advisor",
				Summary:     "Test Advisor",
				Description: "Check Description",
				Category:    "test_category",
				Tiers:       []common.Tier{common.Anonymous},
				Checks: []Check{
					{
						Version:     1,
						Name:        "test_check",
						Summary:     "Test Check",
						Description: "Check Description",
						Advisor:     "test_advisor",
						Type:        MySQLSelect,
						Query:       "id, name FROM table WHERE id=123;",
						Script:      "def func(args): pass",
					},
					{
						Version:     1,
						Name:        "test_check",
						Summary:     "Test Check",
						Description: "Check Description",
						Advisor:     "test_advisor",
						Type:        MySQLSelect,
						Query:       "id, name FROM table WHERE id=123;",
						Script:      "def func(args): pass",
					},
				},
			},
			errStr: "check name collision `test_check` detected in 'test_advisor' advisor",
		},
		{
			name: "normal no tiers specified",
			advisor: &Advisor{
				Version:     1,
				Name:        "test_advisor",
				Summary:     "Test Advisor",
				Description: "Advisor Description",
				Category:    "test_category",
			},
			errStr: "",
		},
		{
			name: "missing name",
			advisor: &Advisor{
				Version:     1,
				Name:        "",
				Summary:     "Test Advisor",
				Description: "Advisor Description",
				Category:    "test_category",
				Tiers:       []common.Tier{common.Anonymous},
			},
			errStr: "invalid advisor name",
		},
		{
			name: "missing summary",
			advisor: &Advisor{
				Version:     1,
				Name:        "test_advisor",
				Summary:     "",
				Description: "Advisor Description",
				Category:    "test_category",
				Tiers:       []common.Tier{common.Anonymous},
			},
			errStr: "summary is empty",
		},
		{
			name: "missing description",
			advisor: &Advisor{
				Version:     1,
				Name:        "test_advisor",
				Summary:     "Test Advisor",
				Description: "",
				Category:    "test_category",
				Tiers:       []common.Tier{common.Anonymous},
			},
			errStr: "description is empty",
		},
		{
			name: "missing category",
			advisor: &Advisor{
				Version:     1,
				Name:        "test_advisor",
				Summary:     "Test Advisor",
				Description: "Advisor Description",
				Category:    "",
				Tiers:       []common.Tier{common.Anonymous},
			},
			errStr: "category is empty",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.advisor.Validate()

			if tt.errStr != "" {
				assert.EqualError(t, err, tt.errStr)
				return
			}

			require.NoError(t, err)
		})
	}
}
