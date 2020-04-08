package checked

import (
	"io"

	"github.com/pkg/errors"
	"go.starlark.net/starlark"
	"gopkg.in/yaml.v3"
)

type Check struct {
	Type   string `yaml:"type"`
	Query  string `yaml:"query"`
	Script string `yaml:"script"`
}

func Parse(r io.Reader) ([]*Check, error) {
	d := yaml.NewDecoder(r)

	type checks struct {
		Checks []*Check `yaml:"checks"`
	}

	res := make([]*Check, 0)
	for {
		c := &checks{}
		err := d.Decode(c)
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

func (c *Check) Validate() error {
	if err := c.validateType(); err != nil {
		return err
	}

	if c.Query == "" {
		return errors.New("check query is empty")
	}

	if err := c.validateScript(); err != nil {
		return err
	}

	return nil
}

func (c *Check) validateScript() error {
	if _, _, err := starlark.SourceProgram("", c.Script, func(s string) bool { return false }); err != nil {
		return errors.Wrap(err, "script is invalid")
	}
	return nil
}

func (c *Check) validateType() error {
	switch c.Type {
	case "MYSQL_SHOW":
		fallthrough
	case "MYSQL_SELECT":
		fallthrough
	case "POSTGRESQL_SHOW":
		fallthrough
	case "POSTGRESQL_SELECT":
		return nil
	case "":
		return errors.New("check type is empty")
	default:
		return errors.Errorf("unknown check type: %s", c.Type)
	}
}

func (c *Check) createFullQuery() string {
	switch c.Type {
	case "MYSQL_SHOW":
		fallthrough
	case "POSTGRESQL_SHOW":
		return "SHOW " + c.Query

	case "MYSQL_SELECT":
		fallthrough
	case "POSTGRESQL_SELECT":
		return "SELECT " + c.Query
	}

	return ""
}
