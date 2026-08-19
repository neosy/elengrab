package dlexecutor

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *Executor) ExecuteRefreshMetadataTask(
	ctx context.Context,
	workerID uint64,
	task *ddownload.RefreshMetadataTask,
) error {
	media, err := uc.download.GetByDownloadIDNoCache(ctx, task.DownloadID)
	if err != nil {
		return uc.failRefreshMetadataTask(ctx, workerID, task, err)
	}

	metadataPatch, err := uc.collectMetadata(ctx, media)
	if err != nil {
		return uc.failRefreshMetadataTask(ctx, workerID, task, err)
	}

	err = uc.applyMetadataPatch(ctx, media, metadataPatch)
	if err != nil {
		return uc.failRefreshMetadataTask(ctx, workerID, task, err)
	}

	err = uc.downloadStatus.Done(ctx, task.DownloadID, nil)
	if err != nil {
		uc.logger.Warn(
			"failed to update download status to done",
			"downloadID", task.DownloadID,
			"error", err,
		)
		return err
	}

	uc.broadcaster.DownloadUpdate(ctx, task.DownloadID)

	uc.logger.Debug(
		"Metadata refresh task completed",
		"workerID", workerID,
		"downloadID", task.DownloadID,
		"title", media.MediaTitle,
	)

	return nil
}

func (uc *Executor) failRefreshMetadataTask(
	ctx context.Context,
	workerID uint64,
	task *ddownload.RefreshMetadataTask,
	err error,
) error {
	updateErr := uc.downloadStatus.Done(ctx, task.DownloadID, nil)
	if updateErr != nil {
		uc.logger.Warn(
			"failed to update download status to done",
			"downloadID", task.DownloadID,
			"error", updateErr,
		)
	}

	uc.logger.Error(
		"metadata refresh task failed",
		"workerID", workerID,
		"downloadID", task.DownloadID,
		"error", err,
	)

	return err
}
