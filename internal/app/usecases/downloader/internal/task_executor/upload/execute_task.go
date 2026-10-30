package uploadexecutor

import (
	"context"
	"errors"
	"strings"

	apperrors "github.com/neosy/elengrab/internal/app/errors"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *Executor) Execute(
	ctx context.Context,
	workerId uint64,
	task *ddownload.MediaProcessingTask,
) error {
	if task == nil {
		uc.logger.Error("Nil pointer in function for upload task", "funcName", "Execute")
		return apperrors.ErrFuncParamNullPointer
	}

	defer func() {
		uc.removeSourceFile(task)
		uc.updateSystemDiskInfo()
		uc.broadcaster.MediaUpdated(ctx, task.DownloadID)
	}()

	err := uc.uploadStatus.Working(ctx, task.DownloadID, task.TaskID, workerId)
	if err != nil {
		uc.logger.Error("Failed update status", "error", err)
		uc.uploadStatus.Failed(ctx, task.DownloadID, nil, err.Error())
		return err
	}

	uc.broadcaster.MediaUpdated(ctx, task.DownloadID)

	download, err := uc.download.GetByDownloadID(ctx, task.DownloadID)
	if err != nil {
		uc.uploadStatus.Failed(ctx, task.DownloadID, nil, err.Error())
		return err
	}

	if download.MediaInfo == nil {
		err := errors.New("download media info is nil")
		uc.logger.Error("Failed to upload: media info is missing",
			"downloadId", task.DownloadID,
			"taskId", task.TaskID,
			"workerId", workerId,
			"error", err,
		)
		uc.uploadStatus.Failed(ctx, task.DownloadID, nil, err.Error())
		return err
	}

	storedFile, err := uc.saveFile(task, download)
	if err != nil {
		uc.uploadStatus.Failed(ctx, task.DownloadID, nil, err.Error())
		return err
	}

	patch := func(dl *ddownload.MediaDownload) error {
		dl.FileNameNoExt = storedFile.fileNameNoExt
		dl.FileExtNoDot = strings.TrimPrefix(storedFile.fileExt, ".")

		dl.FileName = storedFile.fileName
		dl.SafeReadableFileName = dl.NormalizeFileName()

		dl.PartialHash = storedFile.partialHash

		if storedFile.fileSize != 0 {
			dl.FileSize = &storedFile.fileSize
		}

		return nil
	}

	err = uc.uploadStatus.Done(ctx, task.DownloadID, patch)
	if err != nil {
		uc.uploadStatus.Failed(ctx, task.DownloadID, patch, err.Error())
		return err
	}

	if !uc.needPostProcessing(download) {
		return nil
	}

	uc.broadcaster.MediaUpdated(ctx, task.DownloadID)

	err = uc.postProcessing(ctx, task.DownloadID)
	if err != nil {
		uc.logger.Warn(
			"Failed to process file",
			"downloadId", task.DownloadID,
			"fileName", storedFile.fileName,
			"error", err,
		)
	}

	return nil
}
