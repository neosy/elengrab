package dtypes

import (
	"time"

	"github.com/neosy/elengrab/internal/pkg/dbutils"
)

type QueryOptions struct {
	Before *time.Time

	Limit  *uint64
	Offset *uint64

	Filters  QueryFiltersList
	OrderBys dbutils.OrderByList
}

type QueryMediaOptions struct {
	QueryOptions

	Visibility *QueryMediaVisibility
}
