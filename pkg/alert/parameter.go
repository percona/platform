package alert

import "github.com/pkg/errors"

// Parameter represents alerting rule parameter.
type Parameter struct {
	Name  string        `yaml:"name"`  // required
	Help  string        `yaml:"help"`  // required
	Unit  string        `yaml:"unit"`  // required
	Type  Type          `yaml:"type"`  // required
	Range []interface{} `yaml:"range"` // required
	Value interface{}   `yaml:"value"` // required
}

// GetValueForFloat casts parameter value to the float64.
func (p *Parameter) GetValueForFloat() (float64, error) {
	return castValueToFloat64(p.Value)
}

// GetRangeForFloat casts range parameters to the float64.
func (p Parameter) GetRangeForFloat() (float64, float64, error) {
	var lower, higher float64
	var err error

	if lower, err = castValueToFloat64(p.Range[0]); err != nil {
		return 0, 0, err
	}
	if higher, err = castValueToFloat64(p.Range[1]); err != nil {
		return 0, 0, err
	}

	return lower, higher, nil
}

// Validate validates parameter.
func (p *Parameter) Validate() error {
	var err error
	if p.Name == "" {
		return errors.New("parameter name is empty")
	}

	if p.Help == "" {
		return errors.New("parameter help is empty")
	}

	if p.Unit == "" {
		return errors.New("parameter unit is empty")
	}

	if err = p.Type.Validate(); err != nil {
		return err
	}

	if err = p.validateValue(); err != nil {
		return err
	}

	if err = p.validateRange(); err != nil {
		return err
	}

	return nil
}

func (p *Parameter) validateValue() error {
	switch p.Type {
	case Float:
		if _, err := p.GetValueForFloat(); err != nil {
			return err
		}
	default:
		return errors.Errorf("unknown parameter type: %s", p.Type)
	}

	return nil
}

func (p *Parameter) validateRange() error {
	if len(p.Range) != 2 {
		return errors.Errorf("range should have only two elements, but has %d", len(p.Range))
	}

	switch p.Type {
	case Float:
		if _, err := castValueToFloat64(p.Range[0]); err != nil {
			return errors.Wrapf(err, "invalid lower element of range")
		}
		if _, err := castValueToFloat64(p.Range[1]); err != nil {
			return errors.Errorf("invalid higher element of range")
		}

	default:
		return errors.Errorf("unknown parameter type: %s", p.Type)
	}

	return nil
}
