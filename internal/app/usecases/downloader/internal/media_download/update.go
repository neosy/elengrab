package mediadownload

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *MediaDownload) Update(ctx context.Context, userID *uuid.UUID, download *ddownload.MediaDownload) error {
	download.NormalizeForSave()

	if userID != nil && *userID != uuid.Nil && (download.UserID == nil || *download.UserID == uuid.Nil) {
		download.UserID = userID
	}

	err := uc.downloadRepo().Update(ctx, download)
	if err != nil {
		uc.logger.Warn("Update record error", "error", err)
		return err
	}

	uc.updateCacheDependencies(ctx, download)
	uc.updateDependencies(ctx, download)

	return err
}

func (uc *MediaDownload) updateCacheDependencies(ctx context.Context, download *ddownload.MediaDownload) error {
	uc.downloadCacheRep.Delete(ctx, download.DownloadID)

	if download.Status == dtypes.MediaDownloadStatusWorking {
		uc.dlStateCache.SaveByDownload(ctx, download)
	} else {
		uc.dlStateCache.Delete(ctx, download.DownloadID)
	}

	return nil
}

func (uc *MediaDownload) updateDependencies(ctx context.Context, download *ddownload.MediaDownload) error {
	go func() {
		uc.searchIndex.SaveMediaDownload(ctx, download)
	}()
	return nil
}
