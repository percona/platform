// Package fspconv is used to convert gRPC data structures for filtering, sorting and pagination into Go structs
// that are used by fsp package.
package fspconv

import (
	api "github.com/percona-platform/platform/gen/utils/fsp"
	"github.com/percona-platform/platform/pkg/fsp"
)

// NewPageTotals is a constructor of *api.PageTotals for paginated response.
func NewPageTotals(apiModel *api.FilteringSortingPagination, totalItems uint32) *api.PageTotals {
	if apiModel == nil || apiModel.GetPageParams() == nil {
		return nil
	}

	pt := fsp.NewPageTotals(apiModel.GetPageParams().PageSize, totalItems)

	return &api.PageTotals{
		TotalItems: pt.TotalItems,
		TotalPages: pt.TotalPages,
	}
}

// NewFSP is a primary constructor of FilteringSortingPagination struct that is used by storage layer.
func NewFSP(apiModel *api.FilteringSortingPagination, cfg *fsp.Config) (*fsp.FilteringSortingPagination, error) {
	var err error
	out := new(fsp.FilteringSortingPagination)

	if apiModel.GetFilters() != nil {
		out.Filters = newFilters(apiModel.GetFilters())
	}

	if apiModel.GetPageParams() != nil {
		out.PaginationParams = newPageParams(apiModel.GetPageParams())
	}

	if apiModel.GetSortingParams() != nil {
		out.SortingParams, err = newSortingParams(apiModel.GetSortingParams(), cfg.AllowedSortingColumns)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

func newPageParams(apiModel *api.PageParams) *fsp.PaginationParams {
	return fsp.NewPaginationParams(apiModel.GetPageSize(), apiModel.GetIndex())
}

func newSortingParams(apiModel *api.SortingParams, allowedColumns map[string]struct{}) (*fsp.SortingParams, error) {
	return fsp.NewSortingParams(apiModel.GetFieldName(),
		fsp.SortingOrder(apiModel.GetOrder().String()),
		allowedColumns,
	)
}

func newFilter(apiFilter *api.Filter) fsp.Filter {
	return fsp.NewFilter(fsp.FilterType(apiFilter.GetFilterType().String()), newField(apiFilter.GetField()))
}

func newField(apiField *api.Field) fsp.Field {
	return fsp.NewField(apiField.GetName(), apiField.GetValue().AsInterface())
}

func newFilters(apiFilters []*api.Filter) []fsp.Filter {
	out := make([]fsp.Filter, 0, len(apiFilters))
	for _, af := range apiFilters {
		out = append(out, newFilter(af))
	}
	return out
}
