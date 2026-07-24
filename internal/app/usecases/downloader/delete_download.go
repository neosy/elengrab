package downloader

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"time"

	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// DeleteDownload deletes a download from the system.
func (uc *Downloader) DeleteDownload(
	ctx context.Context,
	authCtx dauth.UserContext,
	downloadID uuid.UUID,
) error {
	var (
		needDeleteFileOnStorage bool
		fileFullName            string
	)

	err := uc.validateWriteOperation(authCtx)
	if err != nil {
		return err
	}

	download, err := uc.download.GetByDownloadIDNoCache(ctx, downloadID)
	if err != nil {
		return err
	}

	err = uc.validateDownloadWriteAccess(authCtx, download)
	if err != nil {
		return err
	}

	fnDelete := func(ctx context.Context, download *ddownload.MediaDownload) error {
		err := uc.download.HardDelete(ctx, download.DownloadID)
		if err != nil {
			uc.logger.Error("Failed to delete download", "downloadID", download.DownloadID, "error", err)
			return err
		}
		uc.broadcastDownloadDelete(ctx, download)
		return nil
	}

	fn := func(ctx context.Context) error {
		download, err := uc.download.FindByDownloadIDNoCache(ctx, downloadID)
		if err != nil {
			return err
		}

		// do not return an error if the download is not found
		if download == nil {
			return nil
		}

		switch download.Status {
		case dtypes.MediaDownloadStatusNew:
			if err := fnDelete(ctx, download); err != nil {
				return err
			}
		case dtypes.MediaDownloadStatusPending, dtypes.MediaDownloadStatusWorking:
			task := download.DownloadTask
			if task == nil || task.JobID == nil {
				return errors.New("jobID is missing")
			}
			if !uc.downloadDispatcher.CancelJob(task.JobID.String()) {
				return errors.New("job cannot be cancelled")
			}
			if err := fnDelete(ctx, download); err != nil {
				return err
			}
		case dtypes.MediaDownloadStatusDone, dtypes.MediaDownloadStatusFailed:
			if err := fnDelete(ctx, download); err != nil {
				return err
			}
			needDeleteFileOnStorage = true
		default:
			err := errors.New("the deletion action is not defined")
			uc.logger.Error(err.Error(), "downloadID", downloadID, "status", download.Status, "error", err)
			return err
		}

		fileFullName = download.FileFullName

		return nil
	}

	err = uc.download.Tx(ctx, fn)
	if err != nil {
		return err
	}

	if needDeleteFileOnStorage && fileFullName != "" {
		go func() {
			err := uc.deleteFileWithRetry(ctx, fileFullName, 10, 5*time.Second)
			if err != nil {
				uc.logger.Warn("Failed delete download", "filePath", uc.downloadsStorage.Path(fileFullName), "error", err)
			}
			uc.UpdateSystemInfo()
		}()
	}

	return nil
}

// deleteFileWithRetry attempts to delete a download at the specified path with retries.
func (uc *Downloader) deleteFileWithRetry(ctx context.Context, fileName string, retries int, retryDelay time.Duration) error {
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
