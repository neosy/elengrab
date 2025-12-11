package ytdownloader

import (
	"context"
	"errors"
	"os"
	"path"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *YouTubeDownloader) DeleteDownload(ctx context.Context, fileId uuid.UUID) error {
	var (
		needDeleteFileOnStorage bool
		fileFullName            string
	)

	fnDelete := func(ctx context.Context, file *ddownload.File) error {
		err := uc.file.HardDelete(ctx, fileId)
		if err != nil {
			uc.logger.Error("Failed to delete file", "fileId", fileId, "error", err)
			return err
		}
		uc.dlState.Delete(ctx, file.FileId)
		return nil
	}

	fn := func(ctx context.Context) error {
		file, err := uc.file.FindByFileId(ctx, fileId, false)
		if err != nil {
			return err
		}

		// do not return an error if the file is not found
		if file == nil {
			return nil
		}

		switch file.Status {
		case dtypes.FileStatusNew:
			if err := fnDelete(ctx, file); err != nil {
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
			if err := fnDelete(ctx, file); err != nil {
				return err
			}
		case dtypes.FileStatusDone, dtypes.FileStatusFailed:
			if err := fnDelete(ctx, file); err != nil {
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

	err := uc.file.FileRep.Tx(ctx, fn)
	if err != nil {
		return err
	}

	if needDeleteFileOnStorage && fileFullName != "" {
		fPath := path.Join(uc.downloadsDir, fileFullName)
		if err := os.Remove(fPath); err != nil {
			uc.logger.Warn("Failed delete file", "filePath", fPath, "error", err)
		}
	}

	return nil
}
