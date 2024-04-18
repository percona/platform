package alert

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnit_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		unit   Unit
		errStr string
	}{
		{
			name:   "seconds",
			unit:   Seconds,
			errStr: "",
		},
		{
			name:   "percentage",
			unit:   Percentage,
			errStr: "",
		},
		{
			name:   "empty",
			unit:   Unit(""),
			errStr: "",
		},
		{
			name:   "unknown",
			unit:   Unit("unknown"),
			errStr: "unhandled parameter unit 'unknown'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.unit.Validate()

			if tt.errStr != "" {
				assert.EqualError(t, err, tt.errStr)
				return
			}

			require.NoError(t, err)
		})
	}
}
