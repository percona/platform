// Package check implements checks parsing and validating.
package check

import (
	"io"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// ParseChecks returns slice of checks parsed from YAML passed via reader.
// Can handle multi-document YAMLs, in that case output will be
// union of checks presented in each file.
func ParseChecks(reader io.Reader) ([]*Check, error) {
	d := yaml.NewDecoder(reader)

	type checks struct {
		Checks []*Check `yaml:"checks"`
	}

	var res []*Check
	for {
		var c checks
		err := d.Decode(&c)
		if err != nil {
			if err != io.EOF {
				return nil, errors.Wrap(err, "failed to parse checks")
			}
			break
		}

		res = append(res, c.Checks...)
	}

	return res, nil
}

// ParseCheck returns single check parsed from YAML passed as byte slice.
func ParseCheck(b []byte) (*Check, error) {
	var c Check
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, errors.Wrap(err, "failed to parse check")
	}

	return &c, nil
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

// ParseResults returns slice of results parsed from YAML passed via reader.
// Can handle multi-document YAMLs, in that case output will be
// union of results presented in each file.
func ParseResults(reader io.Reader) ([]*Result, error) {
	d := yaml.NewDecoder(reader)

	type results struct {
		Results []*Result `yaml:"results"`
	}

	var res []*Result
	for {
		var r results
		err := d.Decode(&r)
		if err != nil {
			if err != io.EOF {
				return nil, errors.Wrap(err, "failed to parse results")
			}
			break
		}

		res = append(res, r.Results...)
	}

	return res, nil
}

// ParseResult returns single Result parsed from YAML passed as byte slice.
func ParseResult(b []byte) (*Result, error) {
	var r Result
	if err := yaml.Unmarshal(b, &r); err != nil {
		return nil, errors.Wrap(err, "failed to parse result")
	}

	return &r, nil
}

// Result represents check result that has status and message.
// In case of FAIL status, message should contain reason.
type Result struct {
	Status  string `yaml:"status"`
	Message string `yaml:"message"`
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
