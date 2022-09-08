package fsp

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

const (
	// ASC is for ascending.
	ASC SortingOrder = "ASC"
	// DESC if for descending.
	DESC SortingOrder = "DESC"
)

// SortingOrder is asc/desc sorting.
type SortingOrder string

func (s SortingOrder) String() string {
	return string(s)
}

// SortingParams is used to sort results of an API response.
type SortingParams struct {
	AllowedColumns map[string]struct{}
	FieldName      string
	Order          SortingOrder
}

// NewSortingParams is a constructor.
// ErrSortingNotAllowedByColumn is returned if field is not present in allowedColumns.
func NewSortingParams(field string, order SortingOrder, allowedColumns map[string]struct{}) (*SortingParams, error) {
	out := &SortingParams{
		AllowedColumns: allowedColumns,
		FieldName:      field,
		Order:          order,
	}
	err := out.Validate()
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Attach sorting to SQL query builder.
func (s *SortingParams) Attach(builder sq.SelectBuilder) sq.SelectBuilder {
	// calling validate here to double-check
	// it is expected that validate will be called inside SortingParams constructor
	// to simplify signature of query builder functions
	err := s.Validate()
	if err != nil {
		return builder
	}

	if s.Order == ASC {
		return builder.OrderBy(s.FieldName + " asc")
	} else if s.Order == DESC {
		return builder.OrderByClause(s.FieldName + " desc")
	}
	return builder
}

// Validate checks if sorting params are correctly set
// e.g. column provided from API is allowed by this application.
func (s *SortingParams) Validate() error {
	if s.AllowedColumns == nil {
		return nil
	}
	_, allowed := s.AllowedColumns[s.FieldName]
	if !allowed {
		return fmt.Errorf("%w: %s", ErrSortingNotAllowedByColumn, s.FieldName)
	}
	return nil
}
