package uploadexecutor

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"uuid"

	apperrors "github.com/neosy/elengrab/internal/app/errors"
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/helper"
	"github.com/neosy/elengrab/internal/app/utils/filename"
	dconsts "github.com/neosy/elengrab/internal/domain/consts"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/filex"
)

// storedFile describes a source file that has been moved into the storage.
type storedFile struct {
	fileName      string
	fileNameNoExt string

	fileExt string

	fileSize    int64
	partialHash *string

	thumbnailID uuid.UUID
}

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

	storedFile, err := uc.saveAndProcessFile(ctx, task, download)
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

		if storedFile.thumbnailID != uuid.Nil() {
			dl.MediaInfo.FrameThumbnailID = &storedFile.thumbnailID
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

// removeSourceFile deletes the task source file and logs a warning if it fails.
func (uc *Executor) removeSourceFile(task *ddownload.MediaProcessingTask) {
	if err := os.Remove(task.SourcePath); err != nil && !errors.Is(err, fs.ErrNotExist) {
		uc.logger.Warn("failed to remove source file",
			"downloadID", task.DownloadID,
			"sourcePath", task.SourcePath,
			"error", err,
		)
	}
}

// saveAndProcessFile moves the source file into the storage, collects its metadata
// and extracts a thumbnail for videos.
func (uc *Executor) saveAndProcessFile(
	ctx context.Context,
	task *ddownload.MediaProcessingTask,
	download *ddownload.MediaDownload,
) (storedFile, error) {
	var resultNil storedFile

	fileExt := filepath.Ext(task.SourcePath)

	uniqueName := task.DownloadID.String()

	uniqueFileName := uniqueName + fileExt
	if dconsts.DefaultIncludeTitleInFilename {
		uniqueFileName = filename.BuildUnique(uniqueName, download.MediaTitle, fileExt)
	}

	// Size and hash must be taken before the move, while the source path still exists.
	result := storedFile{
		fileName:      uniqueFileName,
		fileNameNoExt: strings.TrimSuffix(uniqueFileName, fileExt),
		fileExt:       fileExt,

		fileSize:    helper.FileSize(task.SourcePath),
		partialHash: helper.PartialHash(task.SourcePath),
	}

	needMove, err := uc.checkNeedMove(task.SourcePath, uniqueFileName)
	if err != nil {
		return resultNil, err
	}

	if needMove {
		err := uc.storage.Move(task.SourcePath, uniqueFileName)
		if err != nil {
			uc.logger.Warn("Failed to move file", "source", task.SourcePath, "target", uniqueFileName, "error", err)
			return resultNil, err
		}

		uc.logger.Info("File moved", "source", task.SourcePath, "target", uniqueFileName)
	}

	filePath := uc.storage.Path(uniqueFileName)

	if download.IsVideo() {
		if imageData := uc.extractThumbnailFromFile(ctx, filePath); imageData != nil {
			result.thumbnailID, _ = uc.CreateThumbnail(ctx, imageData, dtypes.ThumbnailSourceTypeVideoFrame)
		}
	}

	return result, nil
}

// checkNeedMove reports whether the file must be moved from the source to the storage.
// It returns false if the file is already in the storage,
// and an error if the file is in neither place (lost) or the storage check fails.
func (uc *Executor) checkNeedMove(sourcePath, targetName string) (bool, error) {
	dstExists, err := uc.storage.Exists(targetName)
	if err != nil {
		uc.logger.Warn("Failed to check target file", "target", targetName, "error", err)
		return false, err
	}
	if dstExists {
		uc.logger.Info("File already in storage, move skipped", "source", sourcePath, "target", targetName)
		return false, nil
	}

	srcExists, err := filex.FileExists(sourcePath)
	if err != nil {
		uc.logger.Warn("Failed to check source file", "source", sourcePath, "error", err)
		return false, err
	}
	if !srcExists {
		err := fmt.Errorf("file %q is missing in both source and storage (target %q)", sourcePath, targetName)
		uc.logger.Warn("File lost", "source", sourcePath, "target", targetName, "error", err)
		return false, err
	}

	return true, nil
}
