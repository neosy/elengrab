package uploadexecutor

import (
	"context"
	"uuid"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type postProcessingData struct {
	thumbnailID uuid.UUID
}

func (d *postProcessingData) isZero() bool {
	if d == nil {
		return true
	}

	return d.thumbnailID == uuid.Nil()
}

// needPostProcessing checks if the media download requires post-processing based on its type.
func (uc *Executor) needPostProcessing(download *ddownload.MediaDownload) bool {
	if download == nil {
		return false
	}
	return download.IsVideo()
}

// postProcessing handles additional processing after a file has been uploaded,
// such as generating thumbnails for video files.
// It updates the media download record with any new information obtained during processing.
func (uc *Executor) postProcessing(
	ctx context.Context,
	downloadID uuid.UUID,
) error {
	var processingData postProcessingData

	download, err := uc.download.GetByDownloadIDNoCache(ctx, downloadID)
	if err != nil {
		return err
	}

	err = uc.uploadStatus.Refreshing(ctx, downloadID)
	if err != nil {
		return err
	}

	uc.broadcaster.StartRefreshing(ctx, downloadID)

	filePath := uc.storage.Path(download.FileName)

	if download.IsVideo() {
		imageData := uc.extractThumbnailFromFile(ctx, filePath)
		if imageData != nil {
			processingData.thumbnailID, _ = uc.CreateThumbnail(ctx, imageData, dtypes.ThumbnailSourceTypeVideoFrame)
		}
	}

	var patch func(*ddownload.MediaDownload) error

	if !processingData.isZero() {
		patch = func(dl *ddownload.MediaDownload) error {
			dl.MediaInfo.FrameThumbnailID = &processingData.thumbnailID
			return nil
		}
	}

	err = uc.uploadStatus.Done(ctx, download.DownloadID, patch)
	if err != nil {
		uc.logger.Warn(
			"Failed to update file processing status",
			"downloadId", download.DownloadID,
			"fileName", download.FileName,
			"error", err,
		)
		return err
	}

	return nil
}
