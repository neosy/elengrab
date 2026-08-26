package dlexecutor

import (
	"context"
	"sync"

	apperrors "github.com/neosy/elengrab/internal/app/errors"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *Executor) ExecuteDownloadTask(
	ctx context.Context,
	workerId uint64,
	task *ddownload.DownloadTask,
) error {
	if task == nil {
		uc.logger.Error("Nil pointer in function", "func", "ExecuteDownloadTask")
		return apperrors.ErrFuncParamNullPointer
	}

	err := uc.downloadStatus.Working(ctx, task.DownloadID, task.TaskID, workerId)
	if err != nil {
		uc.logger.Error("Failed update status", "error", err)
		return err
	}

	uc.broadcaster.DownloadUpdate(ctx, task.DownloadID)

	var wg sync.WaitGroup
	defer func() {
		wg.Wait()

		// Broadcast update download info to clients
		uc.broadcaster.DownloadUpdate(ctx, task.DownloadID)
	}()

	wg.Go(func() {
		uc.fetchIcon(ctx, task.MediaUrl)
	})

	resultCh, err := uc.startDownload(ctx, task)
	if err != nil {
		// The context was canceled
		if ctx.Err() != nil {
			ctx = uc.appCtx
		}
		uc.downloadStatus.Failed(ctx, task.DownloadID, nil, err.Error())
		return err
	}

	processed, err := uc.processDownloadResults(ctx, task, resultCh)
	if err != nil {
		var patch func(*ddownload.MediaDownload) error
		if processed != nil {
			patch = func(download *ddownload.MediaDownload) error {
				uc.mappers.MapProcessedDownloadToMediaDownload(download, processed)
				return nil
			}
		}
		// The context was canceled
		if ctx.Err() != nil {
			ctx = uc.appCtx
		}
		uc.downloadStatus.Failed(ctx, task.DownloadID, patch, err.Error())
		return err
	}

	patch := func(download *ddownload.MediaDownload) error {
		uc.mappers.MapProcessedDownloadToMediaDownload(download, processed)
		return nil
	}

	err = uc.downloadStatus.Done(ctx, task.DownloadID, patch)
	if err != nil {
		uc.downloadStatus.Failed(ctx, task.DownloadID, patch, err.Error())
		return err
	}

	return nil
}
