package dtypes

import (
	"uuid"

	"github.com/neosy/elengrab/internal/pkg/dbutils"
)

type QueryOptions struct {
	Limit  *uint64
	Offset *uint64

	Filters  *QueryFilters
	OrderBys dbutils.OrderByList
}

type QueryMediaOptions struct {
	QueryOptions

	ViewMode   QueryMediaViewMode
	LastRecord QueryMediaDownloadCursor

	Visibility     *QueryMediaVisibility
	IsGuestRequest bool

	// OwnerIDForInProgress limits in-progress downloads to the specified owner.
	OwnerIDForInProgress uuid.UUID
}

func NewQueryOptions() QueryOptions {
	return QueryOptions{
		Filters: NewQueryFilters(),
	}
}

func NewQueryMediaOptions() QueryMediaOptions {
	return QueryMediaOptions{
		Filters: NewQueryFilters(),
	}
}
