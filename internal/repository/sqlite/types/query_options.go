package types

import (
	"time"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/dbutils"
)

type QueryOptions struct {
	Before *time.Time

	Limit  *uint64
	Offset *uint64

	Filters  dbutils.FiltersByField
	OrderBys dbutils.OrderByList
}

type QueryMediaOptions struct {
	QueryOptions

	Visibility *dtypes.QueryMediaVisibility
}

func NewQueryOptions() QueryOptions {
	return QueryOptions{
		Filters: make(dbutils.FiltersByField),
	}
}

func NewQueryMediaOptions() QueryMediaOptions {
	return QueryMediaOptions{
		QueryOptions: QueryOptions{
			Filters: make(dbutils.FiltersByField),
		},
	}
}
