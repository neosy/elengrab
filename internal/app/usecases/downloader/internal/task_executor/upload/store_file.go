package uploadexecutor

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/helper"
	"github.com/neosy/elengrab/internal/app/utils/filename"
	dconsts "github.com/neosy/elengrab/internal/domain/consts"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/pkg/filex"
)

// storedFile describes a source file that has been moved into the storage.
type storedFile struct {
	fileName      string
	fileNameNoExt string

	fileExt string

	fileSize    int64
	partialHash *string
}

// saveFile
func (uc *Executor) saveFile(
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
