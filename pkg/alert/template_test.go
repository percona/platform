package alert

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/percona/promconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/percona-platform/platform/pkg/common"
)

func TestTemplate_Parse(t *testing.T) {
	t.Parallel()

	document := strings.TrimSpace(`
---
templates:
  # adopted from https://awesome-prometheus-alerts.grep.to/rules#mysql 1.2
  - version: 1
    name: mysql_too_many_connections
    summary: MySQL connections in use
    tiers: [anonymous, registered]
    expr: |-
      max_over_time(mysql_global_status_threads_connected[5m]) / ignoring (job)
      mysql_global_variables_max_connections
      * 100
      > [[ .threshold ]]
    params:
      - name: threshold
        summary: A percentage from configured maximum
        unit: "%"
        type: float
        range: [0, 100]
        value: 80
      - name: duration
        summary: A duration parameter for testing
        unit: s
        type: float
        value: 5
      - name: boolean
        summary: A boolean parameter for testing
        type: bool
        value: false
      - name: string
        summary: A string parameter for testing
        type: string
        value: foo
    for: 5m
    severity: warning
    labels:
      foo: bar
    annotations:
      summary: "MySQL too many connections (instance {{ $labels.instance }})"
      description: |-
        More than [[ .threshold ]]% of MySQL connections are in use on {{ $labels.instance }}
        VALUE = {{ $value }}
        LABELS: {{ $labels }}
`)

	params := &ParseParams{
		DisallowUnknownFields:    true,
		DisallowInvalidTemplates: true,
	}

	rs, err := Parse(bytes.NewReader([]byte(document)), params)
	require.NoError(t, err)
	assert.Len(t, rs, 1)

	r := rs[0]
	assert.Equal(t, "mysql_too_many_connections", r.Name)
	assert.Equal(t, uint32(1), r.Version)
	assert.Equal(t, []common.Tier{common.Anonymous, common.Registered}, r.Tiers)
	assert.Equal(t, common.Warning, r.Severity)
	assert.Equal(t, "MySQL connections in use", r.Summary)
	assert.Equal(t, "max_over_time(mysql_global_status_threads_connected[5m]) / ignoring (job)\nmysql_global_variables_max_connections\n* 100\n> [[ .threshold ]]", r.Expr)
	assert.Equal(t, map[string]string{"foo": "bar"}, r.Labels)
	assert.Len(t, r.Annotations, 2)
	assert.Equal(t, "MySQL too many connections (instance {{ $labels.instance }})", r.Annotations["summary"])
	assert.Equal(t, "More than [[ .threshold ]]% of MySQL connections are in use on {{ $labels.instance }}\nVALUE = {{ $value }}\nLABELS: {{ $labels }}", r.Annotations["description"])
	assert.Len(t, r.Params, 4)

	param := r.Params[0]
	assert.Equal(t, "threshold", param.Name)
	assert.Equal(t, "A percentage from configured maximum", param.Summary)
	assert.Equal(t, Percentage, param.Unit)
	assert.Equal(t, Float, param.Type)

	assert.NotEmpty(t, param.Range)
	lower, higher, err := param.GetRangeForFloat()
	require.NoError(t, err)
	assert.Equal(t, float64(0), lower)
	assert.Equal(t, float64(100), higher)

	fv, err := param.GetValueForFloat()
	require.NoError(t, err)
	assert.Equal(t, float64(80), fv)

	param = r.Params[1]
	assert.Equal(t, "duration", param.Name)
	assert.Equal(t, "A duration parameter for testing", param.Summary)
	assert.Equal(t, Seconds, param.Unit)
	assert.Equal(t, Float, param.Type)
	assert.Empty(t, param.Range)

	param = r.Params[2]
	assert.Equal(t, "boolean", param.Name)
	assert.Equal(t, "A boolean parameter for testing", param.Summary)
	assert.Empty(t, param.Unit)
	assert.Equal(t, Bool, param.Type)
	assert.Empty(t, param.Range)

	bv, err := param.GetValueForBool()
	require.NoError(t, err)
	assert.Equal(t, false, bv)

	param = r.Params[3]
	assert.Equal(t, "string", param.Name)
	assert.Equal(t, "A string parameter for testing", param.Summary)
	assert.Empty(t, param.Unit)
	assert.Equal(t, String, param.Type)
	assert.Empty(t, param.Range)

	sv, err := param.GetValueForString()
	require.NoError(t, err)
	assert.Equal(t, "foo", sv)
}

func TestTemplate_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		template Template
		errStr   string
	}{{
		"normal",
		Template{
			Name:    "some_name",
			Version: 1,
			Summary: "Some summary message",
			Tiers:   []common.Tier{common.Anonymous},
			Expr:    "some_expression[5m]",
			Params: []Parameter{{
				Name:    "param",
				Summary: "Some expression parameter",
				Unit:    Seconds,
				Type:    "float",
				Range:   []interface{}{10, 90},
				Value:   50,
			}},
			For:         promconfig.Duration(10 * time.Minute),
			Severity:    common.Warning,
			Labels:      map[string]string{"label1": "foo", "label2": "bar"},
			Annotations: map[string]string{"annotation1": "faz", "annotation2": "baz"},
		},
		"",
	}, {
		"no range",
		Template{
			Name:    "some name",
			Version: 1,
			Summary: "Some summary message",
			Tiers:   []common.Tier{common.Anonymous},
			Expr:    "some_expression[5m]",
			Params: []Parameter{{
				Name:    "param",
				Summary: "Some expression parameter",
				Unit:    Seconds,
				Type:    "float",
				Value:   50,
			}},
			For:         promconfig.Duration(10 * time.Minute),
			Severity:    common.Warning,
			Labels:      map[string]string{"label1": "foo", "label2": "bar"},
			Annotations: map[string]string{"annotation1": "faz", "annotation2": "baz"},
		},
		"",
	}, {
		"empty name",
		Template{
			Name:    "",
			Version: 1,
			Summary: "Some summary message",
			Tiers:   []common.Tier{common.Anonymous},
			Expr:    "some_expression[5m]",
			Params: []Parameter{{
				Name:    "param",
				Summary: "Some expression parameter",
				Unit:    Seconds,
				Type:    "float",
				Range:   []interface{}{10, 90},
				Value:   50,
			}},
			For:         promconfig.Duration(10 * time.Minute),
			Severity:    common.Warning,
			Labels:      map[string]string{"label1": "foo", "label2": "bar"},
			Annotations: map[string]string{"annotation1": "faz", "annotation2": "baz"},
		},
		"template name is empty",
	}, {
		"invalid version",
		Template{
			Name:    "some_name",
			Version: 0,
			Summary: "Some summary message",
			Tiers:   []common.Tier{common.Anonymous},
			Expr:    "some_expression[5m]",
			Params: []Parameter{{
				Name:    "param",
				Summary: "Some expression parameter",
				Unit:    Seconds,
				Type:    "float",
				Range:   []interface{}{10, 90},
				Value:   50,
			}},
			For:         promconfig.Duration(10 * time.Minute),
			Severity:    common.Warning,
			Labels:      map[string]string{"label1": "foo", "label2": "bar"},
			Annotations: map[string]string{"annotation1": "faz", "annotation2": "baz"},
		},
		"unexpected version 0",
	}, {
		"empty summary",
		Template{
			Name:    "some_name",
			Version: 1,
			Summary: "",
			Tiers:   []common.Tier{common.Anonymous},
			Expr:    "some_expression[5m]",
			Params: []Parameter{{
				Name:    "param",
				Summary: "Some expression parameter",
				Unit:    Seconds,
				Type:    "float",
				Range:   []interface{}{10, 90},
				Value:   50,
			}},
			For:         promconfig.Duration(10 * time.Minute),
			Severity:    common.Warning,
			Labels:      map[string]string{"label1": "foo", "label2": "bar"},
			Annotations: map[string]string{"annotation1": "faz", "annotation2": "baz"},
		},
		"template summary is empty",
	}, {
		"invalid tier",
		Template{
			Name:    "some_name",
			Version: 1,
			Summary: "Some summary message",
			Tiers:   []common.Tier{common.Anonymous, "invalid"},
			Expr:    "some_expression[5m]",
			Params: []Parameter{{
				Name:    "param",
				Summary: "Some expression parameter",
				Unit:    Seconds,
				Type:    "float",
				Range:   []interface{}{10, 90},
				Value:   50,
			}},
			For:         promconfig.Duration(10 * time.Minute),
			Severity:    common.Warning,
			Labels:      map[string]string{"label1": "foo", "label2": "bar"},
			Annotations: map[string]string{"annotation1": "faz", "annotation2": "baz"},
		},
		"unknown check tier: \"invalid\"",
	}, {
		"duplicate tier",
		Template{
			Name:    "some_name",
			Version: 1,
			Summary: "Some summary message",
			Tiers:   []common.Tier{common.Anonymous, common.Anonymous},
			Expr:    "some_expression[5m]",
			Params: []Parameter{{
				Name:    "param",
				Summary: "Some expression parameter",
				Unit:    Seconds,
				Type:    "float",
				Range:   []interface{}{10, 90},
				Value:   50,
			}},
			For:         promconfig.Duration(10 * time.Minute),
			Severity:    common.Warning,
			Labels:      map[string]string{"label1": "foo", "label2": "bar"},
			Annotations: map[string]string{"annotation1": "faz", "annotation2": "baz"},
		},
		"duplicate tier: \"anonymous\"",
	}, {
		"normal",
		Template{
			Name:    "some_name",
			Version: 1,
			Summary: "Some summary message",
			Tiers:   []common.Tier{common.Anonymous},
			Expr:    "",
			Params: []Parameter{{
				Name:    "param",
				Summary: "Some expression parameter",
				Unit:    Seconds,
				Type:    "float",
				Range:   []interface{}{10, 90},
				Value:   50,
			}},
			For:         promconfig.Duration(10 * time.Minute),
			Severity:    common.Warning,
			Labels:      map[string]string{"label1": "foo", "label2": "bar"},
			Annotations: map[string]string{"annotation1": "faz", "annotation2": "baz"},
		},
		"template expression is empty",
	}, {
		"invalid severity",
		Template{
			Name:    "some_name",
			Version: 1,
			Summary: "Some summary message",
			Tiers:   []common.Tier{common.Anonymous},
			Expr:    "some_expression[5m]",
			Params: []Parameter{{
				Name:    "param",
				Summary: "Some expression parameter",
				Unit:    Seconds,
				Type:    "float",
				Range:   []interface{}{10, 90},
				Value:   50,
			}},
			For:         promconfig.Duration(10 * time.Minute),
			Severity:    common.Severity(256),
			Labels:      map[string]string{"label1": "foo", "label2": "bar"},
			Annotations: map[string]string{"annotation1": "faz", "annotation2": "baz"},
		},
		"unknown severity level: Severity(256)",
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.template.Validate()

			if tt.errStr != "" {
				assert.EqualError(t, err, tt.errStr)
				return
			}

			assert.NoError(t, err)
		})
	}
}
