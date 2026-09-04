package mediadownload

import (
	"context"
	"time"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (uc *MediaDownload) FindByDownloadIDNoCache(
	ctx context.Context,
	downloadID uuid.UUID,
) (*ddownload.MediaDownload, error) {
	download, err := uc.downloadRepo().FindByDownloadID(ctx, downloadID)
	if err != nil {
		uc.logger.Warn("Failed to find record", "error", err)
		return nil, err
	}

	if download != nil {
		uc.downloadCacheRep.Save(ctx, download)
	}

	return download, err
}

// GetByDownloadIDNoCache
// MediaDownload MUST exist — otherwise NOT_FOUND
func (uc *MediaDownload) GetByDownloadIDNoCache(
	ctx context.Context,
	downloadID uuid.UUID,
) (*ddownload.MediaDownload, error) {
	download, err := uc.FindByDownloadIDNoCache(ctx, downloadID)
	if err != nil {
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	if download == nil {
		uc.logger.Warn("MediaDownload not found", "downloadID", downloadID)
		return nil, errorx.New("download not found", exceptionx.NOT_FOUND)
	}

	return download, nil
}

func (uc *MediaDownload) FindByDownloadID(
	ctx context.Context,
	downloadID uuid.UUID,
) (*ddownload.MediaDownload, error) {
	if downloadID == uuid.Nil {
		return nil, nil
	}

	mediaDownload, cacheStatus, _ := uc.downloadCacheRep.FindByDownloadID(ctx, downloadID)
	if mediaDownload != nil {
		return mediaDownload, nil
	}
	if cacheStatus == memsimple.CacheStatusNegativeHit {
		return nil, nil
	}

	mediaDownload, err := uc.FindByDownloadIDNoCache(ctx, downloadID)
	if err != nil {
		return nil, err
	}

	if mediaDownload != nil {
		uc.downloadCacheRep.Save(ctx, mediaDownload)
	} else {
		uc.downloadCacheRep.SaveNegative(ctx, downloadID)
	}

	return mediaDownload, nil
}

func (uc *MediaDownload) GetByDownloadID(
	ctx context.Context,
	downloadID uuid.UUID,
) (*ddownload.MediaDownload, error) {
	download, err := uc.FindByDownloadID(ctx, downloadID)
	if err != nil {
		return nil, err
	}

	if download == nil {
		uc.logger.Warn("MediaDownload not found", "downloadID", downloadID)
		return nil, errorx.New("download not found", exceptionx.NOT_FOUND)
	}

	return download, nil
}

func (uc *MediaDownload) iterateGetAll(ctx context.Context, includeDeleted bool, fn func(*ddownload.MediaDownload) error) error {
	repo := uc.downloadRepo()

	if includeDeleted {
		repo = repo.WithDeleted()
	}

	err := repo.IterateGetAll(ctx, fn)
	if err != nil {
		uc.logger.Warn("Failed to get downloads", "error", err)
		return err
	}

	return nil
}

func (uc *MediaDownload) IterateGetAll(ctx context.Context, fn func(*ddownload.MediaDownload) error) error {
	return uc.iterateGetAll(ctx, false, fn)
}

func (uc *MediaDownload) IterateGetAllWithDeleted(ctx context.Context, fn func(*ddownload.MediaDownload) error) error {
	return uc.iterateGetAll(ctx, true, fn)
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
	filters map[dtypes.QueryFilterName]any,
) ([]*ddownload.MediaDownload, error) {
	repo := uc.downloadRepo()

	if queryOptions != nil {
		repo = repo.WithOptions(*queryOptions)
	}

	if filters != nil {
		repo = repo.WithFilters(filters)
	}

	var downloads []*ddownload.MediaDownload

	err := repo.IterateGetAll(ctx, func(download *ddownload.MediaDownload) error {
		downloads = append(downloads, download)
		return nil
	})
	if err != nil {
		var options dtypes.QueryMediaOptions
		if queryOptions != nil {
			options = *queryOptions
		}

		uc.logger.Warn(
			"Failed to get downloads",
			"queryOptions", options,
			"filters", filters,
			"error", err,
		)

		return nil, err
	}

	return downloads, err
}

func (uc *MediaDownload) GetByStatus(ctx context.Context, status dtypes.MediaDownloadStatus) ([]*ddownload.MediaDownload, error) {
	download, err := uc.downloadRepo().GetByStatus(ctx, status)
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

func (u *MediaDownload) IterateGetByIDs(
	ctx context.Context,
	ids []uuid.UUID,
	fn func(*ddownload.MediaDownload) error,
) error {
	err := u.downloadRepo().IterateGetByIDs(ctx, ids, fn)
	if err != nil {
		u.logger.Warn(
			"Failed to get mediaDownload",
			"ids", ids,
			"error", err,
		)
		return err
	}

	return nil
}

func (u *MediaDownload) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*ddownload.MediaDownload, error) {
	downloads, err := u.downloadRepo().GetByIDs(ctx, ids)
	if err != nil {
		u.logger.Warn(
			"Failed to get mediaDownload",
			"ids", ids,
			"error", err,
		)
		return nil, err
	}

	return downloads, nil
}
