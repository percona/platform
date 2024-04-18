package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSeverity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input     string
		severity  Severity
		canonical string
	}{
		{input: "Emergency", severity: Emergency, canonical: "emergency"},
		{input: "alert", severity: Alert, canonical: "alert"},
		{input: " critical", severity: Critical, canonical: "critical"},
		{input: "Error		", severity: Error, canonical: "error"},
		{input: "WARNING", severity: Warning, canonical: "warning"},
		{input: "notice", severity: Notice, canonical: "notice"},
		{input: "Info", severity: Info, canonical: "info"},
		{input: "Debug", severity: Debug, canonical: "debug"},
		{input: "Unknown", severity: Unknown, canonical: "unknown"},
		{input: "", severity: Unknown, canonical: "unknown"},
		{input: "     ", severity: Unknown, canonical: "unknown"},
		{input: "awesome", severity: Unknown, canonical: "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			actual := ParseSeverity(tt.input)

			assert.Equal(t, tt.severity, actual)
			assert.Equal(t, tt.canonical, actual.String())
		})
	}
}
