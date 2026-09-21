package searchindex

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (u *SearchIndex) GetSourceIndexeByDownloadID(
	ctx context.Context,
	downloadID uuid.UUID,
) (*ddownload.MediaSourceIndex, error) {
	return u.sourceIndex.GetByDownloadID(ctx, downloadID)
}

func (u *SearchIndex) IterateSourceIndexes(ctx context.Context, fn func(*ddownload.MediaSourceIndex) error) error {
	return u.sourceIndex.IterateAll(ctx, fn)
}

func (u *SearchIndex) GetSourceIndexes(
	ctx context.Context,
	queryOptions *dtypes.QueryMediaOptions,
) ([]*ddownload.MediaSourceIndex, error) {
	return u.sourceIndex.GetAll(ctx, queryOptions)
}

func (u *SearchIndex) GetDownloadIDsFromMediaSourceIndex(
	ctx context.Context,
	queryOptions *dtypes.QueryMediaOptions,
) ([]uuid.UUID, error) {
	return u.sourceIndex.GetDownloadIDs(ctx, queryOptions)
}
