package sourceindex

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (uc *MediaSourceIndex) FindByDownloadID(
	ctx context.Context,
	downloadID uuid.UUID,
) (*ddownload.MediaSourceIndex, error) {
	download, err := uc.indexRepo().FindByDownloadID(ctx, downloadID)
	if err != nil {
		uc.logger.Warn("Failed to find record", "error", err)
		return nil, err
	}

	return download, err
}

func (uc *MediaSourceIndex) GetByDownloadID(
	ctx context.Context,
	downloadID uuid.UUID,
) (*ddownload.MediaSourceIndex, error) {
	download, err := uc.FindByDownloadID(ctx, downloadID)
	if err != nil {
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	if download == nil {
		uc.logger.Warn("MediaSourceIndex not found", "downloadID", downloadID)
		return nil, errorx.New("mediaSourceIndex not found", exceptionx.NOT_FOUND)
	}

	return download, nil
}

func (u *MediaSourceIndex) IterateAll(ctx context.Context, fn func(*ddownload.MediaSourceIndex) error) error {
	err := u.indexRepo().IterateAll(ctx, fn)
	if err != nil {
		u.logger.Warn("Failed to get sourceIndex", "error", err)
		return err
	}

	return nil
}

func (u *MediaSourceIndex) GetAll(
	ctx context.Context,
	queryOptions *dtypes.QueryMediaOptions,
) ([]*ddownload.MediaSourceIndex, error) {
	repo := u.indexRepo()

	if queryOptions != nil {
		repo = repo.WithOptions(*queryOptions)
	}

	var indexes []*ddownload.MediaSourceIndex

	err := repo.IterateAll(ctx, func(index *ddownload.MediaSourceIndex) error {
		indexes = append(indexes, index)
		return nil
	})
	if err != nil {
		options := dtypes.NewQueryMediaOptions()
		if queryOptions != nil {
			options = *queryOptions
		}

		u.logger.Warn(
			"Failed to get sourceIndexes",
			"queryOptions", options,
			"filters", queryOptions.Filters,
			"error", err)
		return nil, err
	}

	return indexes, err
}

func (u *MediaSourceIndex) GetDownloadIDs(
	ctx context.Context,
	queryOptions *dtypes.QueryMediaOptions,
) ([]uuid.UUID, error) {
	indexes, err := u.GetAll(ctx, queryOptions)
	if err != nil {
		return nil, err
	}

	if len(indexes) == 0 {
		return nil, nil
	}

	ids := make([]uuid.UUID, 0, len(indexes))

	for _, index := range indexes {
		ids = append(ids, index.DownloadID)
	}

	return ids, nil
}
