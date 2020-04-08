package checked

import (
	"io"

	"github.com/pkg/errors"
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
