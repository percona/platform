package fsp

import "github.com/pkg/errors"

//revive:disable
var (
	ErrSortingNotAllowedByColumn = errors.New("sorting not allowed by this column")
)
