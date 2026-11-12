package types

import (
	"uuid"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/dbutils"
)

type QueryOptions struct {
	Limit  *uint64
	Offset *uint64

	Filters  dbutils.FiltersByField
	OrderBys dbutils.OrderByList
}

type QueryMediaOptions struct {
	QueryOptions

	ViewMode   dtypes.QueryMediaViewMode
	LastRecord dtypes.QueryMediaDownloadCursor

	SearchText     *string
	Visibility     *dtypes.QueryMediaVisibility
	IsGuestRequest bool

	// OwnerIDForInProgress limits in-progress downloads to the specified owner.
	OwnerIDForInProgress uuid.UUID
}

func NewQueryOptions() QueryOptions {
	return QueryOptions{
		Filters: make(dbutils.FiltersByField),
	}
}

func NewQueryMediaOptions() QueryMediaOptions {
	return QueryMediaOptions{
		QueryOptions: NewQueryOptions(),
	}
}
