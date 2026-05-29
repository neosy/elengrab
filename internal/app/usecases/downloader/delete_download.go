package downloader

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/exceptions"
)

// DeleteDownload deletes a download from the system.
func (uc *Downloader) DeleteDownload(
	ctx context.Context,
	userCtx dauth.UserContext,
	downloadID uuid.UUID,
) error {
	var (
		needDeleteFileOnStorage bool
		fileFullName            string
	)

	if uc.demoMode {
		uc.broadcastNotification(
			userCtx.EventKey(),
			dto.BroadcastNotificationModuleResultRow,
			dto.BroadcastNotificationTypeError,
			"Operation not allowed in demo mode",
		)
		return exceptions.DEMO_MODE_RESTRICTION.NewErrorx()
	}

	var accessByUserID *uuid.UUID
	if uc.authz.RestrictDownloadsByUser(userCtx.RoleIDs) {
		accessByUserID = &userCtx.UserID
	}

	download, err := uc.download.GetByDownloadID(ctx, accessByUserID, downloadID)
	if err != nil {
		return err
	}

	fnDelete := func(ctx context.Context, downloadID uuid.UUID) error {
		err := uc.download.HardDelete(ctx, downloadID)
		if err != nil {
			uc.logger.Error("Failed to delete download", "downloadID", downloadID, "error", err)
			return err
		}
		uc.broadcastDownloadDelete(download.UserID, downloadID)
		return nil
	}

	fn := func(ctx context.Context) error {
		download, err := uc.download.FindByDownloadID(ctx, nil, downloadID)
		if err != nil {
			return err
		}

		// do not return an error if the download is not found
		if download == nil {
			return nil
		}

		switch download.Status {
		case dtypes.MediaDownloadStatusNew:
			if err := fnDelete(ctx, download.DownloadID); err != nil {
				return err
			}
		case dtypes.MediaDownloadStatusPending, dtypes.MediaDownloadStatusWorking:
			task := download.DownloadTask
			if task == nil || task.JobID == nil {
				return errors.New("jobID is missing")
			}
			if !uc.dlDispetcher.CancelJob(task.JobID.String()) {
				return errors.New("job cannot be cancelled")
			}
			if err := fnDelete(ctx, download.DownloadID); err != nil {
				return err
			}
		case dtypes.MediaDownloadStatusDone, dtypes.MediaDownloadStatusFailed:
			if err := fnDelete(ctx, download.DownloadID); err != nil {
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

	go func() {
		uc.deleteThumbnails(ctx, download)
	}()

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

func (uc *Downloader) deleteThumbnails(ctx context.Context, download *ddownload.MediaDownload) {
	if download == nil || download.MediaInfo == nil {
		return
	}

	if download.MediaInfo.ThumbnailID != nil {
		uc.thumbnail.Delete(ctx, *download.MediaInfo.ThumbnailID)
	}

	if download.MediaInfo.FrameThumbnailID != nil {
		uc.thumbnail.Delete(ctx, *download.MediaInfo.FrameThumbnailID)
	}
}
