package alert

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParameter_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		parameter Parameter
		errStr    string
	}{
		{
			name: "normal float with range and default value",
			parameter: Parameter{
				Name:    "example",
				Summary: "example parameter",
				Unit:    Seconds,
				Type:    Float,
				Range:   []any{0, 2},
				Value:   1.1,
			},
			errStr: "",
		},
		{
			name: "normal float with range",
			parameter: Parameter{
				Name:    "example",
				Summary: "example parameter",
				Unit:    Seconds,
				Type:    Float,
				Range:   []any{0, 2},
			},
			errStr: "",
		},
		{
			name: "normal float with default value",
			parameter: Parameter{
				Name:    "example",
				Summary: "example parameter",
				Unit:    Seconds,
				Type:    Float,
				Value:   1.1,
			},
			errStr: "",
		},
		{
			name: "normal string",
			parameter: Parameter{
				Name:    "example",
				Summary: "example parameter",
				Type:    String,
				Value:   "test",
			},
			errStr: "",
		},
		{
			name: "normal string without default value",
			parameter: Parameter{
				Name:    "example",
				Summary: "example parameter",
				Type:    String,
			},
			errStr: "",
		},
		{
			name: "normal boolean",
			parameter: Parameter{
				Name:    "example",
				Summary: "example parameter",
				Type:    Bool,
				Value:   true,
			},
			errStr: "",
		},
		{
			name: "normal boolean without default value",
			parameter: Parameter{
				Name:    "example",
				Summary: "example parameter",
				Type:    Bool,
			},
			errStr: "",
		},
		{
			name: "missing name",
			parameter: Parameter{
				Name:    "",
				Summary: "example parameter",
				Unit:    Seconds,
				Type:    Float,
				Range:   []any{0, 2},
				Value:   1.1,
			},
			errStr: "parameter name is empty",
		},
		{
			name: "missing summary",
			parameter: Parameter{
				Name:    "example",
				Summary: "",
				Unit:    Seconds,
				Type:    Float,
				Range:   []any{0, 2},
				Value:   1.1,
			},
			errStr: "parameter summary is empty",
		},
		{
			name: "unknown type",
			parameter: Parameter{
				Name:    "example",
				Summary: "example parameter",
				Unit:    Seconds,
				Type:    Type("unknown"),
				Range:   []any{0, 2},
				Value:   1.1,
			},
			errStr: "unhandled parameter type 'unknown'",
		},
		{
			name: "empty type",
			parameter: Parameter{
				Name:    "example",
				Summary: "example parameter",
				Unit:    Seconds,
				Range:   []any{0, 2},
				Value:   1.1,
			},
			errStr: "unhandled parameter type ''",
		},
		{
			name: "unknown unit",
			parameter: Parameter{
				Name:    "example",
				Summary: "example parameter",
				Unit:    Unit("unknown"),
				Type:    Float,
				Range:   []any{0, 2},
				Value:   1.1,
			},
			errStr: "unhandled parameter unit 'unknown'",
		},
		{
			name: "invalid range",
			parameter: Parameter{
				Name:    "example",
				Summary: "example parameter",
				Unit:    Seconds,
				Type:    Float,
				Range:   []any{0},
				Value:   1.1,
			},
			errStr: "range should be empty or have two elements, but it has 1 elements",
		},
		{
			name: "unit",
			parameter: Parameter{
				Name:    "example",
				Summary: "example parameter",
				Unit:    Seconds,
				Type:    Float,
				Range:   []any{0, 2},
				Value:   1.1,
			},
			errStr: "",
		},
		{
			name: "range is unavailable for string parameters",
			parameter: Parameter{
				Name:    "example",
				Summary: "example parameter",
				Type:    String,
				Range:   []any{"a", "z"},
			},
			errStr: "range should be empty, but it has 2 elements",
		},
		{
			name: "range is unavailable for boolean parameters",
			parameter: Parameter{
				Name:    "example",
				Summary: "example parameter",
				Type:    Bool,
				Range:   []any{"a", "z"},
			},
			errStr: "range should be empty, but it has 2 elements",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.parameter.Validate()

			if tt.errStr != "" {
				assert.EqualError(t, err, tt.errStr)
				return
			}

			require.NoError(t, err)
		})
	}
}
