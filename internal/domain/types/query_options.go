package dtypes

import (
	"time"
)

type QueryOptions struct {
	Before          *time.Time
	Limit           *uint64
	Visibility *QueryMediaVisibility
	IsGuestRequest  bool
}

type QueryMediaOptions struct {
	QueryOptions

	Visibility *QueryMediaVisibility
}
