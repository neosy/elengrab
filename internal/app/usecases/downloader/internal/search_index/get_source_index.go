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
	filters map[dtypes.QueryFilterName]any,
) ([]*ddownload.MediaSourceIndex, error) {
	return u.searchIndex.GetAll(ctx, queryOptions, filters)
}

func (u *SearchIndex) GetDownloadIDsFromMediaSourceIndex(
	ctx context.Context,
	queryOptions *dtypes.QueryMediaOptions,
	filters map[dtypes.QueryFilterName]any,
) ([]uuid.UUID, error) {
	return u.searchIndex.GetDownloadIDs(ctx, queryOptions, filters)
}
