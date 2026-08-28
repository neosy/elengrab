package auth

import (
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type queryOptions struct {
	dtypes.QueryOptions

	withoutGuest *bool
}
