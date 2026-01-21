package ytdownloader

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

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
		if err := os.Remove(fPath); err != nil {
			uc.logger.Warn("Failed delete file", "filePath", fPath, "error", err)
		}
	}

	return nil
}
