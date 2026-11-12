package downloader

import (
	"context"

	apperrors "github.com/neosy/elengrab/internal/app/errors"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *downloader) ExecuteUploadTask(
	ctx context.Context,
	workerId uint64,
	task *ddownload.MediaProcessingTask,
) error {
	if task == nil {
		return apperrors.ErrFuncParamNullPointer
	}

	err := uc.taskExecutors.Upload.Execute(ctx, workerId, task)
	if err != nil {
		return err
	}

	downloadInfo, _ := uc.findActualDownloadInfo(ctx, task.DownloadID)
	uc.broadcastMediaDownloadAdded(dauth.NewEmptyAuthContext(), downloadInfo)

	return nil
}
