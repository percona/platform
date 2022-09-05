package fsp

import (
	"math"

	sq "github.com/Masterminds/squirrel"
)

type (
	// PaginationParams is used for paginating the response.
	PaginationParams struct {
		PageSize  uint32 // maximum size of the page
		PageIndex uint32 // index of requested page, starting with 0
	}

	// PageTotals is a part of a paginated response.
	PageTotals struct {
		TotalItems uint32 `json:"total_items"`
		TotalPages uint32 `json:"total_pages"`
	}
)

// NewPaginationParams is a constructor.
func NewPaginationParams(pageSize, pageIndex uint32) *PaginationParams {
	return &PaginationParams{
		PageSize:  pageSize,
		PageIndex: pageIndex,
	}
}

// NewPageTotals is a constructor.
func NewPageTotals(pageSize, totalItems uint32) *PageTotals {
	return &PageTotals{
		TotalItems: totalItems,
		TotalPages: countTotalPages(pageSize, totalItems),
	}
}

// Attach limit and offset based on p data.
func (p *PaginationParams) Attach(builder sq.SelectBuilder) sq.SelectBuilder {
	limit, offset := p.PageSize, p.PageIndex*p.PageSize
	return builder.Limit(uint64(limit)).Offset(uint64(offset))
}

func countTotalPages(pageSize, totalItems uint32) uint32 {
	return uint32(math.Ceil(float64(totalItems) / float64(pageSize)))
}
