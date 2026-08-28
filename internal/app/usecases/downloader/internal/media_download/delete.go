package mediadownload

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *MediaDownload) SoftDelete(ctx context.Context, downloadID uuid.UUID) error {
	download, err := uc.GetByDownloadID(ctx, downloadID)
	if err != nil {
		return err
	}

	err = uc.downloadRepo().SoftDelete(ctx, downloadID)
	if err != nil {
		uc.logger.Warn("Failed delete file", "error", err)
		return err
	}

	uc.downloadCacheRep.Delete(ctx, downloadID)
	uc.dlStateCache.Delete(ctx, downloadID)

	err = uc.softDeleteDependencies(ctx, download)
	if err != nil {
		return err
	}

	return nil
}

func (uc *MediaDownload) HardDelete(ctx context.Context, downloadID uuid.UUID) error {
	download, err := uc.GetByDownloadID(ctx, downloadID)
	if err != nil {
		return err
	}

	err = uc.downloadRepo().HardDelete(ctx, downloadID)
	if err != nil {
		uc.logger.Warn("Failed delete media download", "error", err)
		return err
	}

	uc.downloadCacheRep.Delete(ctx, downloadID)

	err = uc.dlStateCache.Delete(ctx, downloadID)
	if err != nil {
		uc.logger.Warn("Failed delete download state cache", "error", err)
	}

	err = uc.hardDeleteDependencies(ctx, download)
	if err != nil {
		return err
	}

	return nil
}

func (uc *MediaDownload) softDeleteDependencies(ctx context.Context, download *ddownload.MediaDownload) error {
	go func() {
		uc.searchIndex.SoftDeleteMediaDownload(ctx, download.DownloadID)
	}()

	return nil
}

func (uc *MediaDownload) hardDeleteDependencies(ctx context.Context, download *ddownload.MediaDownload) error {
	err := uc.mediaWatch.DeleteAllByDownloadID(ctx, download.DownloadID)
	if err != nil {
		uc.logger.Warn(
			"Failed to delete media watch data",
			"downloadID", download.DownloadID,
			"error", err,
		)
		return err
	}

	go func() {
		uc.deleteThumbnails(ctx, download)
		uc.searchIndex.HardDeleteMediaDownload(ctx, download.DownloadID)
	}()

	return nil
}

func (uc *MediaDownload) deleteThumbnails(ctx context.Context, download *ddownload.MediaDownload) error {
	if download == nil || download.MediaInfo == nil {
		return nil
	}

	thumblIDs := make([]uuid.UUID, 0, 2)

	if download.MediaInfo.ThumbnailID != nil {
		thumblIDs = append(thumblIDs, *download.MediaInfo.ThumbnailID)
	}

	if download.MediaInfo.FrameThumbnailID != nil {
		thumblIDs = append(thumblIDs, *download.MediaInfo.FrameThumbnailID)
	}

	if len(thumblIDs) > 0 {
		err := uc.thumbnail.DeleteBatch(ctx, thumblIDs)
		if err != nil {
			uc.logger.Warn(
				"Failed to delete thumbnails",
				"thumbnailIDs", thumblIDs,
				"error", err,
			)
			return err
		}
	}

	return nil
}
