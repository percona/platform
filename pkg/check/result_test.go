package check

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/percona/platform/pkg/common"
)

func TestCheck_ResultValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		result *Result
		errStr string
	}{
		{
			name:   "normal",
			result: &Result{Severity: common.Notice, Summary: "some text", ReadMoreURL: "https://www.percona.com/"},
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
		{
			name:   "invalid_read_more_url",
			result: &Result{Severity: common.Notice, Summary: "some text", ReadMoreURL: "percona.com"},
			errStr: "read_more_url: percona.com is invalid",
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

			require.NoError(t, err)
		})
	}
}
