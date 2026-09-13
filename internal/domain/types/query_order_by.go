package dtypes

import "github.com/neosy/elengrab/internal/pkg/dbutils"

type (
	QueryOrder   = dbutils.OrderDirection
	QueryOrderBy = dbutils.OrderBy
)

const (
	QueryOrderAsc  = dbutils.OrderAscending
	QueryOrderDesc = dbutils.OrderDescending
)

var (
	QuerySortBy = dbutils.SortBy
)
