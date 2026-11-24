package sourceindex

import (
	"context"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (u *MediaSourceIndex) Count(ctx context.Context, queryOptions *dtypes.QueryMediaOptions) (int, error) {
	repo := u.indexRepo()

	if queryOptions != nil {
		repo = repo.WithOptions(*queryOptions)
	}

	return repo.Count(ctx)
}
