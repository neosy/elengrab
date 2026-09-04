package auth

import "github.com/neosy/elengrab/internal/repository/sqlite/types"

type queryOptions struct {
	types.QueryOptions

	withoutGuest *bool
}

func newQueryOptions() queryOptions {
	return queryOptions{
		QueryOptions: types.NewQueryOptions(),
	}
}
