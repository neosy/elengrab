package downloader

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"time"

	"uuid"

	"github.com/neosy/elengrab/internal/app/usecases/types"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// DeleteDownload deletes a download from the system.
func (uc *downloader) DeleteDownload(
	ctx context.Context,
	authCtx dauth.AuthContext,
	downloadID uuid.UUID,
) error {
	if err := uc.validateWriteAccess(authCtx, types.BroadcastNotificationModuleResultRow); err != nil {
		return err
	}

	download, err := uc.download.GetByDownloadIDNoCache(ctx, downloadID)
	if err != nil {
		return err
	}

	err = uc.validateDownloadDeleteAccess(authCtx, download)
	if err != nil {
		return err
	}

	type processDeleteResponse struct {
		needDeleteFileOnStorage bool
		FileName                string
	}

	processDelete := func(ctx context.Context) (processDeleteResponse, error) {
		var response, responseNil processDeleteResponse

		download, err := uc.download.FindByDownloadIDNoCache(ctx, downloadID)
		if err != nil {
			return responseNil, err
		}

		// do not return an error if the download is not found
		if download == nil {
			return responseNil, nil
		}

		switch download.Status {
		case dtypes.MediaDownloadStatusNew:
			// No action required for this status.
		case dtypes.MediaDownloadStatusPending, dtypes.MediaDownloadStatusWorking:
			switch download.SourceType {
			case dtypes.MediaDownloadSourceTypeURL:
				uc.cancelDownloadJob(download)
			case dtypes.MediaDownloadSourceTypeUpload:
				uc.cancelUploadJob(download)
			default:
				err := errors.New("unsupported download source type")
				uc.logger.Error(
					err.Error(),
					"downloadID", download.DownloadID,
					"sourceType", download.SourceType,
					"error", err,
				)
				return responseNil, err
			}
		case dtypes.MediaDownloadStatusDone, dtypes.MediaDownloadStatusFailed:
			// No action required for this status.
		default:
			err := errors.New("the deletion action is not defined")
			uc.logger.Error(
				err.Error(),
				"downloadID", downloadID,
				"status", download.Status,
				"error", err,
			)
			return responseNil, err
		}

		err = uc.download.HardDelete(ctx, download.DownloadID)
		if err != nil {
			uc.logger.Error("Failed to delete download", "downloadID", download.DownloadID, "error", err)
			return responseNil, err
		}

		uc.broadcastMediaDownloadDeleted(ctx, download)

		response.needDeleteFileOnStorage = download.FileName != ""
		response.FileName = download.FileName

		return response, nil
	}

	var (
		needDeleteFileOnStorage bool
		FileName                string
	)

	err = uc.download.Tx(ctx,
		func(ctx context.Context) error {
			resp, err := processDelete(ctx)
			if err != nil {
				return err
			}

			needDeleteFileOnStorage = resp.needDeleteFileOnStorage
			FileName = resp.FileName

			return nil
		},
	)
	if err != nil {
		return err
	}

	if needDeleteFileOnStorage && FileName != "" {
		go func() {
			err := uc.deleteFileWithRetry(ctx, FileName, 10, 5*time.Second)
			if err != nil {
				uc.logger.Warn("Failed delete download", "filePath", uc.downloadsStorage.Path(FileName), "error", err)
			}
			uc.UpdateSystemDiskInfo()
		}()
	}

	return nil
}

func (uc *downloader) cancelDownloadJob(download *ddownload.MediaDownload) {
	task := download.DownloadTask
	if task == nil || task.JobID == nil {
		uc.logger.Warn("task or jobId is missing",
			"taskIsNil", task == nil,
			"sourceType", download.SourceType,
			"downloadID", download.DownloadID,
		)
		return
	}

	jobID := task.JobID.String()
	if !uc.downloadDispatcher.HasJob(jobID) {
		uc.logger.Warn("job not found in dispatcher",
			"jobId", jobID,
			"sourceType", download.SourceType,
			"downloadID", download.DownloadID,
		)
		return
	}

	uc.downloadDispatcher.CancelJob(jobID)
}

func (uc *downloader) cancelUploadJob(download *ddownload.MediaDownload) {
	task := download.ProcessingTask
	if task == nil || task.JobID == nil {
		uc.logger.Warn("task or jobId is missing",
			"taskIsNil", task == nil,
			"sourceType", download.SourceType,
			"downloadID", download.DownloadID,
		)
		return
	}

	jobID := task.JobID.String()
	if !uc.mediaProcessingDispatcher.HasJob(jobID) {
		uc.logger.Warn("job not found in dispatcher",
			"jobId", jobID,
			"sourceType", download.SourceType,
			"downloadID", download.DownloadID,
		)
		return
	}

	uc.mediaProcessingDispatcher.CancelJob(jobID)
}

// deleteFileWithRetry attempts to delete a download at the specified path with retries.
func (uc *downloader) deleteFileWithRetry(ctx context.Context, fileName string, retries int, retryDelay time.Duration) error {
	var err error
	for range retries {
		err = uc.downloadsStorage.Delete(fileName)
		if err == nil {
			return nil
		}
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("context was closed. The download was not deleted.")
		case <-time.After(retryDelay):
		}

	}
	return err
}
