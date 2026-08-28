package dlstate

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *DownloadStateCache) Save(ctx context.Context, state *ddownload.DownloadState) error {
	return uc.stateCacheRep.Save(ctx, state)
}

func (uc *DownloadStateCache) Patch(
	ctx context.Context,
	downloadID uuid.UUID,
	mutate func(*ddownload.DownloadState) error,
) error {
	return uc.Transaction(func(ctx context.Context) error {
		state, err := uc.GetByDownloadID(ctx, downloadID)
		if err != nil {
			return err
		}

		if err := mutate(state); err != nil {
			return err
		}

		uc.stateCacheRep.Save(ctx, state)

		return nil
	})
}

func (uc *DownloadStateCache) PatchDownload(ctx context.Context, req dto.PatchMediaDownloadRequest) error {
	return uc.Transaction(func(ctx context.Context) error {
		state, _ := uc.FindByDownloadID(ctx, req.DownloadID)
		if state != nil && state.Download != nil {
			if req.MediaTitle != nil {
				state.Download.MediaTitle = *req.MediaTitle
			}
			if req.MediaDescription != nil {
				state.Download.MediaDescription = *req.MediaDescription
			}
			if req.Visibility != nil {
				state.Download.Visibility = *req.Visibility
			}
			uc.Save(ctx, state)
		}
		return nil
	})
}

func (uc *DownloadStateCache) SaveByDownload(ctx context.Context, download *ddownload.MediaDownload) error {
	if download == nil {
		return nil
	}

	var taskId *uuid.UUID
	if download.DownloadTask != nil {
		taskId = &download.DownloadTask.TaskID
	}

	state := &ddownload.DownloadState{
		DownloadID: download.DownloadID,
		TaskID:     taskId,
		Download:   download.Copy(),
	}

	oldState, _ := uc.stateCacheRep.FindByDownloadID(ctx, download.DownloadID)
	if oldState != nil {
		state.Progress = oldState.Progress

		if state.Download.MediaInfo == nil && oldState.Download.MediaInfo != nil {
			state.Download.MediaInfo = oldState.Download.MediaInfo
		}
	}

	err := uc.stateCacheRep.Save(ctx, state)
	if err != nil {
		uc.logger.Warn("Failed to save download state cache", "error", err)
		return err
	}

	return nil
}
