package fsp

import sq "github.com/Masterminds/squirrel"

const (
	// Equal is for SQL "=" if single value is provided, SQL "IN" if array of values are provided.
	Equal FilterType = "EQUAL"
	// NotEqual is for SQL "<>" if single value is provided, SQL "NOT IN" if array of values are provided.
	NotEqual FilterType = "NOT_EQUAL"
	// Greater is for SQL >
	Greater FilterType = "GREATER"
	// GreaterOrEqual is for SQL >=
	GreaterOrEqual FilterType = "GREATER_OR_EQUAL"
	// Less is for SQL <
	Less FilterType = "LESS"
	// LessOrEqual is for SQL <=
	LessOrEqual FilterType = "LESS_OR_EQUAL"
)

//nolint:gochecknoglobals
var (
	filterToGenerateSQLizerFunc = map[FilterType]generateSQLizerFunc{
		Equal:          generateEqualsSQLizer,
		NotEqual:       generateNotEqualsSQLizer,
		Greater:        generateGreaterThanSQLizer,
		GreaterOrEqual: generateGreaterOrEqualsThanSQLizer,
		Less:           generateLessThanSQLizer,
		LessOrEqual:    generateLessOrEqualsThanSQLizer,
	}
)

type (
	// Filter is a filter for SQL query.
	Filter struct {
		FilterType FilterType
		Field      Field
	}

	// Field represents a field (name and value) of an object.
	// When used with SQL - this is basically a column name and column value.
	Field struct {
		Name  string
		Value interface{}
	}

	// FilterType specifies what kind of filter is this (equals, gte, etc...)
	FilterType string

	generateSQLizerFunc func(Field) sq.Sqlizer
)

// NewFilter is a constructor.
func NewFilter(ft FilterType, field Field) Filter {
	return Filter{
		FilterType: ft,
		Field:      field,
	}
}

// NewField is a constructor.
func NewField(name string, value interface{}) Field {
	return Field{
		Name:  name,
		Value: value,
	}
}

// Attach f to query builder.
func (f *Filter) Attach(builder sq.SelectBuilder) sq.SelectBuilder {
	return builder.Where(f.generateSQLizer())
}

func (f Filter) generateSQLizer() sq.Sqlizer {
	generateFn, ok := filterToGenerateSQLizerFunc[f.FilterType]
	if !ok {
		return nil
	}
	return generateFn(f.Field)
}

func generateEqualsSQLizer(f Field) sq.Sqlizer {
	return sq.Eq{f.Name: f.Value}
}

func generateNotEqualsSQLizer(f Field) sq.Sqlizer {
	return sq.NotEq{f.Name: f.Value}
}

func generateGreaterThanSQLizer(f Field) sq.Sqlizer {
	return sq.Gt{f.Name: f.Value}
}

func generateGreaterOrEqualsThanSQLizer(f Field) sq.Sqlizer {
	return sq.GtOrEq{f.Name: f.Value}
}

func generateLessThanSQLizer(f Field) sq.Sqlizer {
	return sq.Lt{f.Name: f.Value}
}

func generateLessOrEqualsThanSQLizer(f Field) sq.Sqlizer {
	return sq.LtOrEq{f.Name: f.Value}
}
