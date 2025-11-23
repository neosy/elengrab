package ytdownloader

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *YouTubeDownloader) saveStateByFile(ctx context.Context, file *ddownload.File) {
	if err := uc.dlState.SaveByFile(ctx, file); err != nil {
		uc.logger.Warn("Save record failed", "error", err)
	}
}

func (uc *YouTubeDownloader) saveStateByFileId(ctx context.Context, fileId uuid.UUID) {
	if file, _ := uc.file.FindByFileId(ctx, fileId, false); file != nil {
		uc.saveStateByFile(ctx, file)
	}
}
