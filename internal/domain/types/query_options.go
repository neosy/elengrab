package dtypes

import (
	"time"
)

type QueryOptions struct {
	Before *time.Time

	Limit  *uint64
	Offset *uint64

	OrderBy []QueryOrderBy
}

type QueryMediaOptions struct {
	QueryOptions

	Visibility *QueryMediaVisibility
}
