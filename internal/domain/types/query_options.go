package dtypes

import (
	"github.com/neosy/elengrab/internal/pkg/dbutils"
)

type QueryOptions struct {
	Limit  *uint64
	Offset *uint64

	// CreatedAtFilter *dbutils.DateFilter

	Filters  QueryFiltersList
	OrderBys dbutils.OrderByList
}

type QueryMediaOptions struct {
	QueryOptions

	ViewMode   QueryMediaViewMode
	LastRecord QueryMediaDownloadCursor

	Visibility *QueryMediaVisibility
}
