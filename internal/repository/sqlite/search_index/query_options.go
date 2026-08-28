package searchindex

import (
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type queryOptions struct {
	dtypes.QueryMediaOptions

	includeDeleted bool
}
