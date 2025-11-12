package ytdownloader

import (
	"context"
	"errors"
	"fmt"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

func (uc *YouTubeDownloader) ExecuteDownloadTask(
	ctx context.Context,
	workerId uint,
	task *ddownload.DownloadTask,
) error {
	if task == nil {
		uc.logger.Error("Nil pointer in function")
		return errors.New("function parameter is a null pointer")
	}

	err := uc.fileStatus.Working(ctx, task.FileId, task.TaskId, workerId)
	if err != nil {
		return err
	}

	result, err := uc.downloaderSrv.Download(task.YoutubeUrl, task.Options)
	if err != nil {
		uc.fileStatus.Failed(ctx, task.FileId, uptr.String(err.Error()))
		return err
	}

	safeReadableFullName := fmt.Sprintf("%s.%s", uc.sanitizeFileName(result.YoutubeTitle), result.FileExt)

	patch := &dto.FileInfoPatch{
		YoutubeTitle:         &result.YoutubeTitle,
		FileName:             &result.Filename,
		Ext:                  &result.FileExt,
		FullName:             &result.FileFullName,
		SafeReadableFullName: &safeReadableFullName,
	}

	err = uc.fileStatus.Done(ctx, task.FileId, patch)
	if err != nil {
		uc.fileStatus.Failed(ctx, task.FileId, uptr.String(err.Error()))
		return err
	}

	return nil
}
