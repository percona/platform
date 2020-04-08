package checked

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var singleFileData = `---
checks:
  - type: MYSQL_SHOW
    query: VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');
    script: |
        def function1(args):
            pass
`

var multipleFilesData = `---
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

func TestCheck_Parse(t *testing.T) {

	t.Run("singleFile", func(t *testing.T) {
		cs, err := Parse(bytes.NewReader([]byte(singleFileData)))
		require.NoError(t, err)

		assert.Len(t, cs, 1)
		assert.Equal(t, "MYSQL_SHOW", cs[0].Type)
		assert.Equal(t, "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');", cs[0].Query)
		assert.Equal(t, cs[0].Script, "def function1(args):\n    pass\n")
	})


	t.Run("multipleFiles", func(t *testing.T) {
		cs, err := Parse(bytes.NewReader([]byte(multipleFilesData)))
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

func TestCheck_validateType(t *testing.T) {
	tests := []struct {
		name   string
		typ    string
		errStr string
	}{
		{name: "mysql_show", typ: "MYSQL_SHOW", errStr: ""},
		{name: "mysql_select", typ: "MYSQL_SELECT", errStr: ""},
		{name: "postgresql_show", typ: "POSTGRESQL_SHOW", errStr: ""},
		{name: "postgresql_select", typ: "POSTGRESQL_SHOW", errStr: ""},
		{name: "clickhouse_show", typ: "CLICKHOUSE_SHOW", errStr: "unknown check type: CLICKHOUSE_SHOW"},
		{name: "empty", typ: "", errStr: "check type is empty"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Check{Type: tt.typ}
			err := c.validateType()

			if tt.errStr != "" {
				assert.EqualError(t, err, tt.errStr)
				return
			}

			assert.NoError(t, err)
		})
	}
}
