package mediadownload

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
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

	return download, err
}

func (uc *MediaDownload) ExistsByDownloadIDNoCache(
	ctx context.Context,
	downloadID uuid.UUID,
) (bool, error) {
	exists, err := uc.downloadRepo().ExistsByDownloadID(ctx, downloadID)
	if err != nil {
		uc.logger.Warn("Failed to find record", "error", err)
		return false, err
	}

	return exists, err
}

// GetByDownloadIDNoCache
// MediaDownload MUST exist — otherwise NOT_FOUND
func (uc *MediaDownload) GetByDownloadIDNoCache(
	ctx context.Context,
	downloadID uuid.UUID,
) (*ddownload.MediaDownload, error) {
	if downloadID == uuid.Nil {
		uc.logger.Warn("MediaDownload not found", "downloadID", downloadID)
		return nil, nil
	}

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

	if mediaDownload == nil {
		uc.downloadCacheRep.SaveNegative(ctx, downloadID, idcodec.EncodeUUIDShortBase64URL(downloadID))
		return nil, nil
	}

	uc.downloadCacheRep.Save(ctx, mediaDownload)

	return mediaDownload, nil
}

func (uc *MediaDownload) GetByDownloadID(
	ctx context.Context,
	downloadID uuid.UUID,
) (*ddownload.MediaDownload, error) {
	if downloadID == uuid.Nil {
		uc.logger.Warn("MediaDownload not found", "downloadID", downloadID)
		return nil, nil
	}

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

func (u *MediaDownload) IterateByIDs(
	ctx context.Context,
	ids []uuid.UUID,
	fn func(*ddownload.MediaDownload) error,
) error {
	err := u.downloadRepo().IterateByIDs(ctx, ids, fn)
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

func (u *MediaDownload) GetByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) ([]*ddownload.MediaDownload, error) {
	repo := u.downloadRepo()

	downloads, err := repo.GetByIDs(ctx, ids)
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
