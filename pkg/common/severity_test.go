package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSeverity(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		severity Severity
	}{
		{name: "normal", str: "Emergency", severity: Emergency},
		{name: "first lowercase ", str: "alert", severity: Alert},
		{name: "first space", str: " critical", severity: Critical},
		{name: "last tabs", str: "Error		", severity: Error},
		{name: "all capital", str: "WARNING", severity: Warning},
		{name: "normal", str: "notice", severity: Notice},
		{name: "normal", str: "Info", severity: Info},
		{name: "normal", str: "Debug", severity: Debug},
		{name: "normal", str: "Unknown", severity: Unknown},
		{name: "empty string", str: "", severity: Unknown},
		{name: "spaces", str: "     ", severity: Unknown},
		{name: "unknown", str: "awesome", severity: Unknown},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := ParseSeverity(tt.str)

			assert.Equal(t, tt.severity, actual)
		})
	}
}
