package searchindex

import (
	"context"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (u *SearchIndex) CountFromIndex(
	ctx context.Context,
	queryOptions *dtypes.QueryMediaOptions,
) (int, error) {
	count, err := u.sourceIndex.Count(ctx, queryOptions)
	if err != nil {
		u.logger.Warn("Failed to count media from source index", "error", err)
	}

	return count, nil
}
