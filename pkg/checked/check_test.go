package checked

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var data = `---
checks:
  - type: MYSQL_SHOW
    query: VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');
    script: |
        def helper(args):
            pass
 
        def check(rows):
            vars = {
                "have_ssl":     "YES",
                "have_openssl": "YES",
            }
 
            for row in rows:
                name = row["Variable_name"]
                actual = row["Value"]
                expected = vars.get(name)
                if expected and expected != actual:
                    return {"error": "expected %s to be %s, got %s" % (name, expected, actual)}
 
            return {}`

func TestCheck_Parse(t *testing.T) {
	cs, err := Parse(bytes.NewReader([]byte(data)))
	require.NoError(t, err)

	require.NoError(t, err)

	assert.Len(t, cs, 1)
	assert.Equal(t, "MYSQL_SHOW", cs[0].Type)
	assert.Equal(t, "VARIABLES WHERE Variable_name IN ('have_ssl', 'have_openssl');", cs[0].Query)

	expectedScript := `def helper(args):
    pass

def check(rows):
    vars = {
        "have_ssl":     "YES",
        "have_openssl": "YES",
    }

    for row in rows:
        name = row["Variable_name"]
        actual = row["Value"]
        expected = vars.get(name)
        if expected and expected != actual:
            return {"error": "expected %s to be %s, got %s" % (name, expected, actual)}

    return {}`

	assert.Equal(t, cs[0].Script, expectedScript)
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
