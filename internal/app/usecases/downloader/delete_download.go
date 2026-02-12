package downloader

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// DeleteDownload deletes a download from the system.
func (uc *YouTubeDownloader) DeleteDownload(
	ctx context.Context,
	userID uuid.UUID,
	fileId uuid.UUID,
) error {
	var (
		needDeleteFileOnStorage bool
		fileFullName            string
	)

	var accessByUserID *uuid.UUID
	if uc.historyMode != dtypes.HistoryModeGlobal {
		accessByUserID = &userID
	}

	_, err := uc.file.GetByFileId(ctx, accessByUserID, fileId)
	if err != nil {
		return err
	}

	fnDelete := func(ctx context.Context, fileID uuid.UUID) error {
		err := uc.file.HardDelete(ctx, fileID)
		if err != nil {
			uc.logger.Error("Failed to delete file", "fileId", fileID, "error", err)
			return err
		}
		return nil
	}

	fn := func(ctx context.Context) error {
		file, err := uc.file.FindByFileId(ctx, nil, fileId)
		if err != nil {
			return err
		}

		// do not return an error if the file is not found
		if file == nil {
			return nil
		}

		switch file.Status {
		case dtypes.FileStatusNew:
			if err := fnDelete(ctx, file.FileId); err != nil {
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
			if err := fnDelete(ctx, file.FileId); err != nil {
				return err
			}
		case dtypes.FileStatusDone, dtypes.FileStatusFailed:
			if err := fnDelete(ctx, file.FileId); err != nil {
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

	if needDeleteFileOnStorage && fileFullName != "" {
		fPath := filepath.Join(uc.downloadsDir, fileFullName)
		go func(path string) {
			err := uc.deleteWithRetry(ctx, path, 10, 5*time.Second)
			if err != nil {
				uc.logger.Warn("Failed delete file", "filePath", fPath, "error", err)
			}
		}(fPath)
	}

	return nil
}

// deleteWithRetry attempts to delete a file at the specified path with retries.
func (uc *YouTubeDownloader) deleteWithRetry(ctx context.Context, path string, retries int, retryDelay time.Duration) error {
	var err error
	for range retries {
		err = os.Remove(path)
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
