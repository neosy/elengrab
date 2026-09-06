package searchindex

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (u *SearchIndex) IterateGetSourceIndexes(ctx context.Context, fn func(*ddownload.MediaSourceIndex) error) error {
	return u.searchIndex.IterateGetAll(ctx, fn)
}

func (u *SearchIndex) GetSourceIndexes(
	ctx context.Context,
	queryOptions *dtypes.QueryMediaOptions,
) ([]*ddownload.MediaSourceIndex, error) {
	return u.searchIndex.GetAll(ctx, queryOptions)
}

func (u *SearchIndex) GetDownloadIDsFromMediaSourceIndex(
	ctx context.Context,
	queryOptions *dtypes.QueryMediaOptions,
) ([]uuid.UUID, error) {
	return u.searchIndex.GetDownloadIDs(ctx, queryOptions)
}
