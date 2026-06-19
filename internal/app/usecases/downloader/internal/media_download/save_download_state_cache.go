package mediadownload

import (
	"context"

	"github.com/google/uuid"
)

func (uc *MediaDownload) saveToDownloadStateCache(ctx context.Context, downloadID uuid.UUID) {
	download, _ := uc.GetByDownloadID(ctx, downloadID)
	if download != nil {
		err := uc.dlStateCache.SaveByDownload(ctx, download)
		if err != nil {
			uc.logger.Warn("Failed to save download state cache", "error", err)
		}
	}
}
