package mediadownload

import (
	"context"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *MediaDownload) DeleteBroken(ctx context.Context) error {
	files, err := uc.GetByStatus(ctx, dtypes.MediaDownloadStatusNew)
	if err != nil {
		uc.logger.Warn("Failed to get files", "error", err)
		return err
	}

	for _, file := range files {
		if file.DownloadTask == nil {
			err := uc.HardDelete(ctx, file.DownloadID)
			if err != nil {
				uc.logger.Warn("Failed to delete file", "downloadID", file.DownloadID, "error", err)
				continue
			}
		}
	}

	return nil
}
