package watchevent

import (
	"time"

	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

type queryOptions struct {
	beforeTime *time.Time
	limit      *uint64
}

func (o *queryOptions) copy() queryOptions {
	if o == nil {
		return queryOptions{}
	}

	options := *o

	options.beforeTime = uptr.Copy(o.beforeTime)
	options.limit = uptr.Copy(o.limit)

	return options
}
