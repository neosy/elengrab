package ytdownloader

import (
	"context"
	"errors"
	"fmt"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/pkg/nfile"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

func (uc *YouTubeDownloader) ExecuteDownloadTask(
	ctx context.Context,
	workerId uint,
	task *ddownload.DownloadTask,
) error {
	if task == nil {
		uc.logger.Error("Nil pointer in function", "func", "ExecuteDownloadTask")
		return errors.New("function parameter is a null pointer")
	}

	err := uc.fileStatus.Working(ctx, task.FileId, task.TaskId, workerId)
	if err != nil {
		uc.logger.Error("Failed update status", "error", err)
		return err
	}

	uc.saveStateByFileId(ctx, task.FileId)

	resultCh, err := uc.downloaderSrv.Download(task.YoutubeUrl, uc.mappers.MapDownloadOptionsDomainToService(task.Options))
	if err != nil {
		uc.logger.Error(
			"Download",
			"error", err,
		)
		uc.fileStatus.Failed(ctx, task.FileId, nil, uptr.String(err.Error()))
		uc.saveStateByFileId(ctx, task.FileId)
		return err
	}

	var lastResult *ddownload.DownloadResult
	for r := range resultCh {
		if r.Error != nil {
			uc.logger.Error(
				"Download",
				"error", r.Error,
			)
			var patch *dto.FileInfoPatch
			if lastResult.YoutubeTitle != "" {
				patch = &dto.FileInfoPatch{
					YoutubeTitle: &lastResult.YoutubeTitle,
				}
				if lastResult.FileExt != "" {
					patch.Ext = &lastResult.FileExt
				}
				if lastResult.Filesize != nil && *lastResult.Filesize != 0 {
					patch.FileSize = &lastResult.Filesize
				}
			}
			uc.fileStatus.Failed(ctx, task.FileId, patch, uptr.String(r.Error.Error()))
			uc.saveStateByFileId(ctx, task.FileId)
			return r.Error
		}

		state, err := uc.dlState.FindByFileId(ctx, task.FileId)
		if err != nil {
			uc.logger.Error(
				"Find by fileId",
				"error", err,
			)
			uc.fileStatus.Failed(ctx, task.FileId, nil, uptr.String(err.Error()))
			uc.saveStateByFileId(ctx, task.FileId)
			return err
		}

		lastResult = r

		state.InitFromDownloadResult(r)
		uc.dlState.Save(
			ctx,
			state,
		)
	}

	safeReadableFullName := fmt.Sprintf("%s.%s", nfile.SanitizeFileName(lastResult.YoutubeTitle), lastResult.FileExt)

	patch := &dto.FileInfoPatch{
		YoutubeTitle:         &lastResult.YoutubeTitle,
		FileName:             &lastResult.Filename,
		Ext:                  &lastResult.FileExt,
		FullName:             &lastResult.FileFullName,
		FileSize:             &lastResult.Filesize,
		PartialHash:          &lastResult.PartialHash,
		SafeReadableFullName: &safeReadableFullName,
		MediaInfo:            &lastResult.MediaInfo,
	}

	err = uc.fileStatus.Done(ctx, task.FileId, patch)
	if err != nil {
		uc.fileStatus.Failed(ctx, task.FileId, patch, uptr.String(err.Error()))
		uc.saveStateByFileId(ctx, task.FileId)
		return err
	}

	uc.saveStateByFileId(ctx, task.FileId)

	return nil
}
