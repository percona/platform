package check

import (
	"net/url"

	"github.com/pkg/errors"

	"github.com/percona-platform/platform/pkg/common"
)

// Result represents a single check script result that is used to generate alert.
type Result struct {
	Summary     string            `json:"summary"`      // required
	Description string            `json:"description"`  // optional
	ReadmoreURL string            `yaml:"readmore_url"` // optional
	Severity    common.Severity   `json:"severity"`     // required
	Labels      map[string]string `json:"labels"`       // optional
}

// Validate validates check result for minimal correctness.
func (r *Result) Validate() error {
	if r.Summary == "" {
		return errors.New("summary is empty")
	}

	if r.ReadmoreURL != "" {
		_, err := url.ParseRequestURI(r.ReadmoreURL)
		if err != nil {
			return errors.Errorf("readmore_url: %s is invalid", r.ReadmoreURL)
		}
	}

	if err := r.Severity.Validate(); err != nil {
		return err
	}

	if r.Severity < common.Error || r.Severity > common.Notice {
		// until UI is ready to support more severities
		return errors.Errorf("unhandled result severity: %s", r.Severity)
	}

	return nil
}
