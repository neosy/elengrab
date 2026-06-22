package downloader

import (
	"context"
	"time"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

const (
	// failedDownloadRetentionPeriod defines how long failed download records
	// are kept in the database before being deleted.
	failedDownloadRetentionPeriod = 24 * time.Hour
)

// DeleteFailedDownloads removes database records for videos
// that were not successfully downloaded from YouTube.
func (uc *Downloader) DeleteFailedDownloads(ctx context.Context) error {
	downloads, err := uc.download.GetByStatus(ctx, dtypes.MediaDownloadStatusFailed)
	if err != nil {
		uc.logger.Error("Failed to get by status", "status", dtypes.MediaDownloadStatusFailed, "error", err)
		return err

	}

	for _, download := range downloads {
		if time.Now().UTC().After(download.CreatedAt.UTC().Add(failedDownloadRetentionPeriod)) {
			err := uc.download.HardDelete(ctx, download.DownloadID)
			if err != nil {
				uc.logger.Warn("Failed to hard delete download", "downloadID", download.DownloadID, "error", err)
				continue
			}
			uc.deleteThumbnails(ctx, download)
			uc.logger.Debug("Hard deleted download download", "downloadID", download.DownloadID)
			uc.broadcastDownloadDelete(ctx, download)
		}
	}

	return nil
}
