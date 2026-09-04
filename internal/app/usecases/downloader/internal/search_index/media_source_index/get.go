package sourceindex

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (uc *MediaSourceIndex) FindByDownload(
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

func (uc *MediaSourceIndex) GetByDownload(
	ctx context.Context,
	downloadID uuid.UUID,
) (*ddownload.MediaSourceIndex, error) {
	download, err := uc.FindByDownload(ctx, downloadID)
	if err != nil {
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	if download == nil {
		uc.logger.Warn("MediaSourceIndex not found", "downloadID", downloadID)
		return nil, errorx.New("mediaSourceIndex not found", exceptionx.NOT_FOUND)
	}

	return download, nil
}

func (u *MediaSourceIndex) IterateGetAll(ctx context.Context, fn func(*ddownload.MediaSourceIndex) error) error {
	err := u.indexRepo().IterateGetAll(ctx, fn)
	if err != nil {
		u.logger.Warn("Failed to get sourceIndex", "error", err)
		return err
	}

	return nil
}

func (u *MediaSourceIndex) GetAll(
	ctx context.Context,
	queryOptions *dtypes.QueryMediaOptions,
	filters map[dtypes.QueryFilterName]any,
) ([]*ddownload.MediaSourceIndex, error) {
	repo := u.indexRepo()

	if queryOptions != nil {
		repo = repo.WithOptions(*queryOptions)
	}

	if len(filters) > 0 {
		repo = repo.WithFilters(filters)
	}

	var indexes []*ddownload.MediaSourceIndex

	err := repo.IterateGetAll(ctx, func(index *ddownload.MediaSourceIndex) error {
		indexes = append(indexes, index)
		return nil
	})
	if err != nil {
		var options dtypes.QueryMediaOptions
		if queryOptions != nil {
			options = *queryOptions
		}

		u.logger.Warn(
			"Failed to get sourceIndexes",
			"queryOptions", options,
			"filters", filters,
			"error", err)
		return nil, err
	}

	return indexes, err
}

func (u *MediaSourceIndex) GetDownloadIDs(
	ctx context.Context,
	queryOptions *dtypes.QueryMediaOptions,
	filters map[dtypes.QueryFilterName]any,
) ([]uuid.UUID, error) {
	indexes, err := u.GetAll(ctx, queryOptions, filters)
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
