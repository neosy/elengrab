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

// DeleteFile deletes a download from the system.
func (uc *Downloader) DeleteFile(
	ctx context.Context,
	userCtx dauth.UserContext,
	fileId uuid.UUID,
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
	if uc.authz.RestrictFilesByUser(userCtx.Roles) {
		accessByUserID = &userCtx.UserID
	}

	file, err := uc.file.GetByFileID(ctx, accessByUserID, fileId)
	if err != nil {
		return err
	}

	fnDelete := func(ctx context.Context, fileID uuid.UUID) error {
		err := uc.file.HardDelete(ctx, fileID)
		if err != nil {
			uc.logger.Error("Failed to delete file", "fileId", fileID, "error", err)
			return err
		}
		uc.broadcastFileDelete(file.UserID, fileID)
		return nil
	}

	fn := func(ctx context.Context) error {
		file, err := uc.file.FindByFileID(ctx, nil, fileId)
		if err != nil {
			return err
		}

		// do not return an error if the file is not found
		if file == nil {
			return nil
		}

		switch file.Status {
		case dtypes.FileStatusNew:
			if err := fnDelete(ctx, file.FileID); err != nil {
				return err
			}
		case dtypes.FileStatusPending, dtypes.FileStatusWorking:
			task := file.DownloadTask
			if task == nil || task.JobID == nil {
				return errors.New("jobID is missing")
			}
			if !uc.dlDispetcher.CancelJob(task.JobID.String()) {
				return errors.New("job cannot be cancelled")
			}
			if err := fnDelete(ctx, file.FileID); err != nil {
				return err
			}
		case dtypes.FileStatusDone, dtypes.FileStatusFailed:
			if err := fnDelete(ctx, file.FileID); err != nil {
				return err
			}
			needDeleteFileOnStorage = true
		default:
			err := errors.New("the deletion action is not defined")
			uc.logger.Error(err.Error(), "fileId", fileId, "status", file.Status, "error", err)
			return err
		}

		fileFullName = file.FullName

		return nil
	}

	err = uc.file.Tx(ctx, fn)
	if err != nil {
		return err
	}

	go func() {
		uc.deleteThumbnails(ctx, file)
	}()

	if needDeleteFileOnStorage && fileFullName != "" {
		go func() {
			err := uc.deleteFileWithRetry(ctx, fileFullName, 10, 5*time.Second)
			if err != nil {
				uc.logger.Warn("Failed delete file", "filePath", uc.downloadsStorage.Path(fileFullName), "error", err)
			}
			uc.UpdateSystemInfo()
		}()
	}

	return nil
}

// deleteFileWithRetry attempts to delete a file at the specified path with retries.
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
			return fmt.Errorf("context was closed. The file was not deleted.")
		case <-time.After(retryDelay):
		}

	}
	return err
}

func (uc *Downloader) deleteThumbnails(ctx context.Context, file *ddownload.File) {
	if file == nil || file.MediaInfo == nil {
		return
	}

	if file.MediaInfo.ThumbnailID != nil {
		uc.thumbnail.Delete(ctx, *file.MediaInfo.ThumbnailID)
	}

	if file.MediaInfo.FrameThumbnailID != nil {
		uc.thumbnail.Delete(ctx, *file.MediaInfo.FrameThumbnailID)
	}
}
