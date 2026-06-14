package dtypes

import "time"

type QueryOptions struct {
	Before          *time.Time
	Limit           *uint64
	MediaVisibility *QueryMediaVisibility
}
