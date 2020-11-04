package check

import (
	"github.com/pkg/errors"

	"github.com/percona-platform/platform/pkg/alert"
)

// Result represents a single check script result that is used to generate alert.
type Result struct {
	Summary     string            `json:"summary"`     // required
	Description string            `json:"description"` // optional
	Severity    alert.Severity    `json:"severity"`    // required
	Labels      map[string]string `json:"labels"`      // optional
}

// Validate validates check result for minimal correctness.
func (r *Result) Validate() error {
	if r.Summary == "" {
		return errors.New("summary is empty")
	}

	if r.Severity < alert.Emergency || r.Severity > alert.Debug {
		return errors.Errorf("unknown result severity: %s", r.Severity)
	}

	if r.Severity < alert.Error || r.Severity > alert.Notice {
		// until UI is ready to support more severities
		return errors.Errorf("unhandled result severity: %s", r.Severity)
	}

	return nil
}
