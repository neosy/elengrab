package fileuc

import (
	"context"

	"github.com/google/uuid"
)

func (uc *File) saveToDownloadStateCache(ctx context.Context, fileId uuid.UUID) {
	file, _ := uc.GetByFileID(ctx, nil, fileId)
	if file != nil {
		err := uc.dlStateCache.SaveByFile(ctx, file)
		if err != nil {
			uc.logger.Warn("Failed to save download state cache", "error", err)
		}
	}
}
