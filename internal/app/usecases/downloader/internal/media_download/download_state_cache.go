package mediadownload

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *MediaDownload) saveToDownloadStateCache(ctx context.Context, downloadID uuid.UUID) {
	download, _ := uc.GetByDownloadIDNoCache(ctx, downloadID)
	if download != nil {
		err := uc.dlStateCache.SaveByDownload(ctx, download)
		if err != nil {
			uc.logger.Warn("Failed to save download state cache", "error", err)
		}
	}
}

func (uc *MediaDownload) FindState(
	ctx context.Context,
	downloadID uuid.UUID,
) (*ddownload.DownloadState, error) {
	return uc.dlStateCache.FindByDownloadID(ctx, downloadID)
}

func (uc *MediaDownload) GetOrCreateState(
	ctx context.Context,
	downloadID uuid.UUID,
) (*ddownload.DownloadState, error) {
	state, err := uc.dlStateCache.FindByDownloadID(ctx, downloadID)
	if err != nil {
		return nil, err
	}

	if state != nil {
		return state, nil
	}

	download, err := uc.GetByDownloadID(ctx, downloadID)
	if err != nil {
		return nil, err
	}

	uc.dlStateCache.SaveByDownload(ctx, download)

	state, err = uc.dlStateCache.GetByDownloadID(ctx, downloadID)
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (uc *MediaDownload) SaveState(
	ctx context.Context,
	state *ddownload.DownloadState,
) error {
	return uc.dlStateCache.Save(ctx, state)
}

func (uc *MediaDownload) PatchState(
	ctx context.Context,
	downloadID uuid.UUID,
	mutate func(*ddownload.DownloadState) error,
) error {
	return uc.dlStateCache.Patch(ctx, downloadID, mutate)
}
