package alert

import "strings"

//go:generate ../../bin/stringer -type=Severity -linecomment

// Severity represents alert severity level.
type Severity int

// Supported severity levels.
const (
	Unknown   Severity = iota // unknown
	Emergency                 // emergency
	Alert                     // alert
	Critical                  // critical
	Error                     // error
	Warning                   // warning
	Notice                    // notice
	Info                      // info
	Debug                     // debug
)

// ParseSeverity casts string to Severity.
func ParseSeverity(s string) Severity {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "emergency":
		return Emergency
	case "alert":
		return Alert
	case "critical":
		return Critical
	case "error":
		return Error
	case "warning":
		return Warning
	case "notice":
		return Notice
	case "info":
		return Info
	case "debug":
		return Debug
	default:
		return Unknown
	}
}
