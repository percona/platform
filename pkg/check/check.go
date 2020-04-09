// Package check implements checks parsing and validating.
package check

import (
	"io"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Parse returns slice of checks parsed from YAML passed via reader.
// Can handle multi-document YAMLs, parsing result will be a single slice
// that contains checks form every parsed document.
func Parse(reader io.Reader) ([]Check, error) {
	d := yaml.NewDecoder(reader)

	type checks struct {
		Checks []Check `yaml:"checks"`
	}

	var res []Check
	for {
		var c checks
		if err := d.Decode(&c); err != nil {
			if err == io.EOF {
				return res, nil
			}
			return nil, errors.Wrap(err, "failed to parse checks")
		}

		res = append(res, c.Checks...)
	}
}

// Supported check types.
const (
	MySQLShow           = "MYSQL_SHOW"
	MySQLSelect         = "MYSQL_SELECT"
	PostgreSQLShow      = "POSTGRESQL_SHOW"
	PostgreSQLSelect    = "POSTGRESQL_SELECT"
	MongoDBGetParameter = "MONGODB_GETPARAMETER"
)

// Check represents security check structure.
type Check struct {
	Type   string `yaml:"type"`
	Query  string `yaml:"query"`
	Script string `yaml:"script"`
}

// Validate validates check for minimal correctness.
func (c *Check) Validate() error {
	if err := c.validateType(); err != nil {
		return err
	}

	if c.Query == "" {
		return errors.New("check query is empty")
	}

	if c.Script == "" {
		return errors.New("check script is empty")
	}

	return nil
}

// validateType validates check type.
func (c *Check) validateType() error {
	switch c.Type {
	case MySQLShow:
		fallthrough
	case MySQLSelect:
		fallthrough
	case PostgreSQLShow:
		fallthrough
	case PostgreSQLSelect:
		fallthrough
	case MongoDBGetParameter:
		return nil
	case "":
		return errors.New("check type is empty")
	default:
		return errors.Errorf("unknown check type: %s", c.Type)
	}
}

// Possible result statuses.
const (
	Success = "SUCCESS"
	Fail    = "FAIL"
)

// Result represents check result that has status and message.
// In case of FAIL status, message should contain reason.
type Result struct {
	Status  string
	Message string
}

// Validate validates check result for minimal correctness.
func (r *Result) Validate() error {
	if err := r.validateStatus(); err != nil {
		return err
	}

	if r.Status == Fail && r.Message == "" {
		return errors.New("failed check result should have message")
	}

	return nil
}

// validateType validates check result status.
func (r *Result) validateStatus() error {
	switch r.Status {
	case Success:
		fallthrough
	case Fail:
		return nil
	case "":
		return errors.New("result status is empty")
	default:
		return errors.Errorf("unknown result status: %s", r.Status)
	}
}
