package dtypes

import "github.com/neosy/elengrab/internal/pkg/dbutils"

type (
	QueryOrder   = dbutils.OrderDirection
	QueryOrderBy = dbutils.OrderBy
)

const (
	QueryOrderAsc  dbutils.OrderDirection = dbutils.OrderAscending
	QueryOrderDesc dbutils.OrderDirection = dbutils.OrderDescending
)
