package dlstate

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *DownloadStateCache) FindByFileId(ctx context.Context, fileId uuid.UUID) (*ddownload.DownloadState, error) {
	return uc.stateRep.FindByFileId(ctx, fileId)
}

func (uc *DownloadStateCache) FindByTaskId(ctx context.Context, taskId uuid.UUID) (*ddownload.DownloadState, error) {
	return uc.stateRep.FindByTaskId(ctx, taskId)
}
