package uploadexecutor

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	apperrors "github.com/neosy/elengrab/internal/app/errors"
	"github.com/neosy/elengrab/internal/app/utils/filename"
	dconsts "github.com/neosy/elengrab/internal/domain/consts"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *Executor) Execute(
	ctx context.Context,
	workerId uint64,
	task *ddownload.MediaProcessingTask,
) error {
	if task == nil {
		return apperrors.ErrFuncParamNullPointer
	}

	err := uc.uploadStatus.Working(ctx, task.DownloadID, task.TaskID, workerId)
	if err != nil {
		uc.logger.Error("Failed update status", "error", err)
		return err
	}

	uc.broadcaster.MediaUpdated(ctx, task.DownloadID)

	download, err := uc.download.GetByDownloadID(ctx, task.DownloadID)
	if err != nil {
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
		return err
	}

	defer func() {
		uc.broadcaster.MediaUpdated(ctx, task.DownloadID)
		uc.updateSystemDiskInfo()
	}()

	uniqueName := task.DownloadID.String()
	var uniqueFileName string

	if dconsts.DefaultIncludeTitleInFilename {
		uniqueFileName = filename.BuildUnique(uniqueName, download.MediaTitle, filepath.Ext(task.SourcePath))
	} else {
		uniqueFileName = uniqueName + filepath.Ext(task.SourcePath)
	}

	fileSize := uc.fileSize(task.SourcePath)
	partialHash := uc.partialHash(task.SourcePath)

	err = uc.storage.Move(task.SourcePath, uniqueFileName)
	if err != nil {
		uc.uploadStatus.Failed(ctx, task.DownloadID, nil, err.Error())
		return err
	}

	filePath := uc.storage.Path(uniqueFileName)

	var thumbnailID uuid.UUID

	if download.IsVideo() {
		imageData := uc.extractThumbnailFromFile(ctx, filePath)
		if imageData != nil {
			thumbnailID, _ = uc.CreateThumbnail(ctx, imageData, dtypes.ThumbnailSourceTypeVideoFrame)
		}
	}

	fileExt := filepath.Ext(uniqueFileName)
	fileName := strings.TrimSuffix(uniqueFileName, fileExt)

	patch := func(dl *ddownload.MediaDownload) error {
		dl.FileName = fileName
		dl.Ext = strings.TrimPrefix(fileExt, ".")

		dl.FileFullName = uniqueFileName
		dl.SafeReadableFileFullName = dl.NormalizeFileFullName()

		dl.PartialHash = partialHash

		if fileSize != 0 {
			dl.FileSize = &fileSize
		}

		if thumbnailID != uuid.Nil {
			dl.MediaInfo.FrameThumbnailID = &thumbnailID
		}

		return nil
	}

	err = uc.uploadStatus.Done(ctx, task.DownloadID, patch)
	if err != nil {
		uc.uploadStatus.Failed(ctx, task.DownloadID, patch, err.Error())
		return err
	}

	return nil
}
