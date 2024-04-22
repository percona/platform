package alert

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/percona/promconfig"
	"github.com/stretchr/testify/require"

	"github.com/percona/platform/pkg/common"
)

func TestTemplate_Parse(t *testing.T) {
	t.Parallel()

	document := strings.TrimSpace(`
templates:
    - name: mysql_too_many_connections
      version: 1
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
          unit: '%'
          type: float
          range: [0, 100]
          value: 80
        - name: duration
          summary: A duration parameter for testing
          unit: s
          type: float
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
        description: |-
            More than [[ .threshold ]]% of MySQL connections are in use on {{ $labels.instance }}
            VALUE = {{ $value }}
            LABELS: {{ $labels }}
        summary: MySQL too many connections (instance {{ $labels.instance }})`)

	params := &ParseParams{
		DisallowUnknownFields:    true,
		DisallowInvalidTemplates: true,
	}

	rs, err := Parse(bytes.NewReader([]byte(document)), params)
	require.NoError(t, err)
	require.Len(t, rs, 1)

	r := rs[0]
	require.Equal(t, "mysql_too_many_connections", r.Name)
	require.Equal(t, uint32(1), r.Version)
	require.Equal(t, []common.Tier{common.Anonymous, common.Registered}, r.Tiers)
	require.Equal(t, common.Warning, r.Severity)
	require.Equal(t, "MySQL connections in use", r.Summary)
	require.Equal(t, "max_over_time(mysql_global_status_threads_connected[5m]) / ignoring (job)\nmysql_global_variables_max_connections\n* 100\n> [[ .threshold ]]", r.Expr)
	require.Equal(t, map[string]string{"foo": "bar"}, r.Labels)
	require.Len(t, r.Annotations, 2)
	require.Equal(t, "MySQL too many connections (instance {{ $labels.instance }})", r.Annotations["summary"])
	require.Equal(t, "More than [[ .threshold ]]% of MySQL connections are in use on {{ $labels.instance }}\nVALUE = {{ $value }}\nLABELS: {{ $labels }}", r.Annotations["description"])
	require.Len(t, r.Params, 4)

	param := r.Params[0]
	require.Equal(t, "threshold", param.Name)
	require.Equal(t, "A percentage from configured maximum", param.Summary)
	require.Equal(t, Percentage, param.Unit)
	require.Equal(t, Float, param.Type)

	require.NotEmpty(t, param.Range)
	lower, higher, err := param.GetRangeForFloat()
	require.NoError(t, err)
	require.InDelta(t, float64(0), lower, 0)
	require.InDelta(t, float64(100), higher, 0)
	fv, err := param.GetValueForFloat()
	require.NoError(t, err)
	require.InDelta(t, float64(80), fv, 0)

	param = r.Params[1]
	require.Equal(t, "duration", param.Name)
	require.Equal(t, "A duration parameter for testing", param.Summary)
	require.Equal(t, Seconds, param.Unit)
	require.Equal(t, Float, param.Type)
	require.Empty(t, param.Range)
	require.Nil(t, param.Value)

	param = r.Params[2]
	require.Equal(t, "boolean", param.Name)
	require.Equal(t, "A boolean parameter for testing", param.Summary)
	require.Empty(t, param.Unit)
	require.Equal(t, Bool, param.Type)
	require.Empty(t, param.Range)
	bv, err := param.GetValueForBool()
	require.NoError(t, err)
	require.False(t, bv)

	param = r.Params[3]
	require.Equal(t, "string", param.Name)
	require.Equal(t, "A string parameter for testing", param.Summary)
	require.Empty(t, param.Unit)
	require.Equal(t, String, param.Type)
	require.Empty(t, param.Range)
	sv, err := param.GetValueForString()
	require.NoError(t, err)
	require.Equal(t, "foo", sv)

	nd, err := ToYAML(rs)
	require.NoError(t, err)
	require.Equal(t, strings.TrimSpace(document), strings.TrimSpace(nd))
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
		"empty expression",
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
		"missing parameter name",
		Template{
			Name:    "some_name",
			Version: 1,
			Summary: "Some summary message",
			Tiers:   []common.Tier{common.Anonymous},
			Expr:    "some_expression[5m]",
			Params: []Parameter{{
				Name:    "",
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
		"parameter '' is invalid: parameter name is empty",
	}, {
		"invalid parameter type",
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
				Type:    "unknown",
				Range:   []interface{}{10, 90},
				Value:   50,
			}},
			For:         promconfig.Duration(10 * time.Minute),
			Severity:    common.Warning,
			Labels:      map[string]string{"label1": "foo", "label2": "bar"},
			Annotations: map[string]string{"annotation1": "faz", "annotation2": "baz"},
		},
		"parameter 'param' is invalid: unhandled parameter type 'unknown'",
	}, {
		"parameter summary is empty",
		Template{
			Name:    "some_name",
			Version: 1,
			Summary: "Some summary message",
			Tiers:   []common.Tier{common.Anonymous},
			Expr:    "some_expression[5m]",
			Params: []Parameter{{
				Name:    "param",
				Summary: "",
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
		"parameter 'param' is invalid: parameter summary is empty",
	}, {
		"missing parameter type",
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
				Type:    "",
				Range:   []interface{}{10, 90},
				Value:   50,
			}},
			For:         promconfig.Duration(10 * time.Minute),
			Severity:    common.Warning,
			Labels:      map[string]string{"label1": "foo", "label2": "bar"},
			Annotations: map[string]string{"annotation1": "faz", "annotation2": "baz"},
		},
		"parameter 'param' is invalid: unhandled parameter type ''",
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.template.Validate()

			if tt.errStr != "" {
				require.EqualError(t, err, tt.errStr)
				return
			}

			require.NoError(t, err)
		})
	}
}
