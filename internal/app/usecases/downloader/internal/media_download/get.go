package mediadownload

import (
	"context"
	"time"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/dbutils"
)

func (uc *MediaDownload) iterateAll(ctx context.Context, includeDeleted bool, fn func(*ddownload.MediaDownload) error) error {
	repo := uc.downloadRepo()

	if includeDeleted {
		repo = repo.WithDeleted()
	}

	err := repo.IterateAll(ctx, fn)
	if err != nil {
		uc.logger.Warn("Failed to get downloads", "error", err)
		return err
	}

	return nil
}

func (uc *MediaDownload) IterateAll(ctx context.Context, fn func(*ddownload.MediaDownload) error) error {
	return uc.iterateAll(ctx, false, fn)
}

func (uc *MediaDownload) IterateAllWithDeleted(ctx context.Context, fn func(*ddownload.MediaDownload) error) error {
	return uc.iterateAll(ctx, true, fn)
}

func (uc *MediaDownload) GetAllFullNames(ctx context.Context, includeDeleted bool) (map[string]struct{}, error) {
	names, err := uc.downloadRepo().GetAllFullNames(ctx, includeDeleted)
	if err != nil {
		uc.logger.Warn("Failed to get fullNames", "error", err)
		return nil, err
	}

	return names, nil
}

func (uc *MediaDownload) GetAll(
	ctx context.Context,
	queryOptions *dtypes.QueryMediaOptions,
) ([]*ddownload.MediaDownload, error) {
	repo := uc.downloadRepo()

	if queryOptions != nil {
		repo = repo.WithOptions(*queryOptions)
	}

	var downloads []*ddownload.MediaDownload

	err := repo.IterateAll(ctx, func(download *ddownload.MediaDownload) error {
		downloads = append(downloads, download)
		return nil
	})
	if err != nil {
		options := dtypes.NewQueryMediaOptions()
		if queryOptions != nil {
			options = *queryOptions
		}

		uc.logger.Warn(
			"Failed to get downloads",
			"queryOptions", options,
			"error", err,
		)

		return nil, err
	}

	return downloads, err
}

func (uc *MediaDownload) GetByStatus(ctx context.Context, status dtypes.MediaDownloadStatus) ([]*ddownload.MediaDownload, error) {
	repo := uc.downloadRepo()
	repo = repo.WithOrderBy(dtypes.QuerySortBy(dtypes.QueryFilterNameCreatedAt.String(), dtypes.QueryOrderAsc))

	download, err := repo.GetByStatus(ctx, status)
	if err != nil {
		uc.logger.Warn("Failed to get downloads", "error", err)
		return nil, err
	}

	return download, err
}

func (uc *MediaDownload) GetByPartialHash(ctx context.Context, criteria ddownload.DuplicateHashRow) ([]*ddownload.MediaDownload, error) {
	downloadRep := uc.downloadRepo()
	if criteria.UserID != nil {
		downloadRep = downloadRep.WithUser(*criteria.UserID)
	}
	download, err := downloadRep.GetByPartialHash(ctx, criteria.Hash)
	if err != nil {
		uc.logger.Warn("Failed to get downloads", "error", err)
		return nil, err
	}

	return download, err
}

func (uc *MediaDownload) GetWithoutPartialHash(ctx context.Context) ([]*ddownload.MediaDownload, error) {
	var downloads []*ddownload.MediaDownload

	gFiles, err := uc.downloadRepo().GetWithoutPartialHash(ctx)
	if err != nil {
		uc.logger.Warn("Failed to get downloads", "error", err)
		return nil, err
	}

	if len(gFiles) > 0 {
		downloads = make([]*ddownload.MediaDownload, 0, len(gFiles))
		for _, download := range gFiles {
			if download.FileFullName == "" {
				continue
			}
			downloads = append(downloads, download)
		}
	}

	return downloads[:len(downloads):len(downloads)], nil
}

func (uc *MediaDownload) GetDuplicateHashes(ctx context.Context, scope dtypes.UniquenessScope) ([]ddownload.DuplicateHashRow, error) {
	rows, err := uc.downloadRepo().GetDuplicateHashes(ctx, scope)
	if err != nil {
		uc.logger.Warn("Failed to get dublicate hashes", "error", err)
		return nil, err
	}

	return rows, nil
}

func (uc *MediaDownload) GetDeleted(ctx context.Context, from, to *time.Time) ([]*ddownload.MediaDownload, error) {
	downloads, err := uc.downloadRepo().GetDeleted(ctx, from, to)
	if err != nil {
		uc.logger.Warn("Failed to get deleted", "fromDate", from, "toDate", to, "error", err)
		return nil, err
	}

	return downloads, nil
}

func (u *MediaDownload) FindLastByChannelID(
	ctx context.Context,
	channelID uuid.UUID,
) (*ddownload.MediaDownload, error) {
	repo := u.downloadRepo()

	repo = repo.WithFilters(
		dtypes.NewQueryFilter(
			dtypes.QueryFilterNameChannelID,
			dbutils.NewFilterConditionEq(channelID),
		),
	)

	repo = repo.WithOrderBy(dtypes.QuerySortBy(dtypes.QueryFilterNameCreatedAt.String(), dbutils.OrderDescending))

	var download *ddownload.MediaDownload

	err := repo.IterateAll(ctx, func(d *ddownload.MediaDownload) error {
		download = d
		return nil
	})
	if err != nil {
		return nil, err
	}

	return download, nil
}
