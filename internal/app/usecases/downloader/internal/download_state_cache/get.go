package dlstate

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *DownloadStateCache) FindByDownloadID(
	ctx context.Context,
	downloadID uuid.UUID,
) (*ddownload.DownloadState, error) {
	return uc.stateRep.FindByDownloadID(ctx, downloadID)
}

func (uc *DownloadStateCache) FindByTaskID(ctx context.Context, taskId uuid.UUID) (*ddownload.DownloadState, error) {
	return uc.stateRep.FindByTaskID(ctx, taskId)
}
