package mediadownload

import (
	"context"

	"uuid"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *MediaDownload) SoftDelete(
	ctx context.Context,
	downloadID uuid.UUID,
) error {
	return uc.SoftDeleteWithOnSuccess(ctx, downloadID, nil)
}

func (uc *MediaDownload) SoftDeleteWithOnSuccess(
	ctx context.Context,
	downloadID uuid.UUID,
	onSuccess func(),
) error {
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

	switch download.SourceType {
	case dtypes.MediaDownloadSourceTypeURL:
		uc.dlStateCache.Delete(ctx, downloadID)
	case dtypes.MediaDownloadSourceTypeUpload:
		uc.uploadStateCacheRepo.Delete(ctx, downloadID)
	}

	if onSuccess != nil {
		onSuccess()
	}

	err = uc.softDeleteDependencies(ctx, download)
	if err != nil {
		return err
	}

	return nil
}

func (uc *MediaDownload) HardDelete(
	ctx context.Context,
	downloadID uuid.UUID,
) error {
	return uc.HardDeleteWithOnSuccess(ctx, downloadID, nil)
}

func (uc *MediaDownload) HardDeleteWithOnSuccess(
	ctx context.Context,
	downloadID uuid.UUID,
	onSuccess func(),
) error {
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

	switch download.SourceType {
	case dtypes.MediaDownloadSourceTypeURL:
		err = uc.dlStateCache.Delete(ctx, downloadID)
		if err != nil {
			uc.logger.Warn("Failed delete download state cache", "error", err)
		}
	case dtypes.MediaDownloadSourceTypeUpload:
		err = uc.uploadStateCacheRepo.Delete(ctx, downloadID)
		if err != nil {
			uc.logger.Warn("Failed delete upload state cache", "error", err)
		}
	}

	if onSuccess != nil {
		onSuccess()
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
		uc.searchIndex.HardDeleteMediaDownload(ctx, download.DownloadID)
		uc.deleteThumbnails(ctx, download)
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
