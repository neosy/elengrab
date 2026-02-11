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
func (uc *YouTubeDownloader) DeleteFailedDownloads(ctx context.Context) error {
	files, err := uc.file.GetByStatus(ctx, dtypes.FileStatusFailed)
	if err != nil {
		uc.logger.Error("Failed to get by status", "status", dtypes.FileStatusFailed, "error", err)
		return err

	}

	for _, file := range files {
		if time.Now().UTC().After(file.CreatedAt.UTC().Add(failedDownloadRetentionPeriod)) {
			err := uc.file.HardDelete(ctx, file.FileId)
			if err != nil {
				uc.logger.Warn("Failed to hard delete file", "fileId", file.FileId, "error", err)
				continue
			}
			uc.logger.Debug("Hard deleted file download", "fileId", file.FileId)
		}
	}

	return nil
}
