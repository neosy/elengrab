package dlstate

import (
	"context"
	"errors"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *DownloadStateCache) FindByDownloadID(
	ctx context.Context,
	downloadID uuid.UUID,
) (*ddownload.DownloadState, error) {
	state, err := uc.stateCacheRep.FindByDownloadID(ctx, downloadID)
	if err != nil {
		uc.logger.Error(
			"Failed to find download state",
			"downloadID", downloadID,
			"error", err,
		)
		return nil, err
	}
	return state, nil
}

func (uc *DownloadStateCache) GetByDownloadID(
	ctx context.Context,
	downloadID uuid.UUID,
) (*ddownload.DownloadState, error) {
	state, err := uc.FindByDownloadID(ctx, downloadID)
	if err != nil {
		return nil, err
	}

	if state == nil {
		return nil, errors.New("download state not found")
	}

	return state, nil
}

func (uc *DownloadStateCache) FindByTaskID(ctx context.Context, taskId uuid.UUID) (*ddownload.DownloadState, error) {
	return uc.stateCacheRep.FindByTaskID(ctx, taskId)
}
