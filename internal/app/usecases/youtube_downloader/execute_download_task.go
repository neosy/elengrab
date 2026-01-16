package ytdownloader

import (
	"context"
	"errors"
	"fmt"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dyoutube "github.com/neosy/elengrab/internal/domain/youtube_info"
	"github.com/neosy/elengrab/pkg/nfasthttp"
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

	resultCh, err := uc.downloaderSrv.Download(ctx, task.YoutubeUrl, uc.mappers.MapDownloadOptionsDomainToService(task.Options))
	if err != nil {
		// The context was canceled
		if ctx.Err() != nil {
			uc.logger.Debug(
				"Failed to download: The context was canceled",
				"error", err,
			)
			file, e := uc.file.FindByFileId(uc.appCtx, nil, task.FileId)
			if e == nil && file != nil {
				uc.fileStatus.Failed(uc.appCtx, task.FileId, nil, uptr.String(err.Error()))
			}
			uc.dlStateCache.Delete(uc.appCtx, task.FileId)
			return ctx.Err()
		}

		uc.logger.Error(
			"Failed to download",
			"error", err,
		)

		uc.fileStatus.Failed(ctx, task.FileId, nil, uptr.String(err.Error()))

		return err
	}

	var lastResult *ddownload.DownloadResult
	for r := range resultCh {
		if r.Error != nil {
			// The context was canceled
			if ctx.Err() != nil {
				file, e := uc.file.FindByFileId(uc.appCtx, nil, task.FileId)
				if e == nil && file != nil {
					uc.fileStatus.Failed(uc.appCtx, task.FileId, nil, uptr.String(r.Error.Error()))
				}
				uc.dlStateCache.Delete(uc.appCtx, task.FileId)
				return ctx.Err()
			}

			uc.logger.Error(
				"Failed to download",
				"fileId", task.FileId,
				"error", r.Error,
			)

			var patch *dto.FileInfoPatch
			if lastResult != nil {
				patch = &dto.FileInfoPatch{
					YoutubeChannelID: &lastResult.ChannelID,
				}
				if lastResult.YoutubeTitle != "" {
					patch.YoutubeTitle = &lastResult.YoutubeTitle
				}
				if lastResult.FileExt != "" {
					patch.Ext = &lastResult.FileExt
				}
				if lastResult.Filesize != nil && *lastResult.Filesize != 0 {
					patch.FileSize = &lastResult.Filesize
				}
			}
			uc.fileStatus.Failed(ctx, task.FileId, patch, uptr.String(r.Error.Error()))
			return r.Error
		}

		state, err := uc.dlStateCache.FindByFileId(ctx, nil, task.FileId)
		if err != nil {
			uc.logger.Error(
				"Failed to download",
				"action", "Find by fileId",
				"error", err,
			)
			uc.fileStatus.Failed(ctx, task.FileId, nil, uptr.String(err.Error()))
			return err
		}

		lastResult = r

		// Adding a record to the YouTube Channel table
		if lastResult != nil && lastResult.ChannelID != nil && lastResult.ChannelAvatar != nil {
			exists, _ := uc.ytChannel.ExistsByChannelID(ctx, *lastResult.ChannelID)
			if !exists {
				channel := &dyoutube.YoutubeChannel{
					ChannelID: *lastResult.ChannelID,
				}
				channel.InitFromResultChannelAvatar(lastResult.ChannelAvatar)
				uc.ytChannel.Create(ctx, channel)
			}
		}

		state.InitFromDownloadResult(r)
		uc.dlStateCache.Save(
			ctx,
			state,
		)
	}

	safeReadableFullName := uptr.String(
		fmt.Sprintf("%s.%s", nfasthttp.SanitizeFileName(lastResult.YoutubeTitle), lastResult.FileExt),
	)

	patch := &dto.FileInfoPatch{
		YoutubeChannelID:     &lastResult.ChannelID,
		YoutubeTitle:         &lastResult.YoutubeTitle,
		FileName:             &lastResult.Filename,
		Ext:                  &lastResult.FileExt,
		FullName:             &lastResult.FileFullName,
		FileSize:             &lastResult.Filesize,
		PartialHash:          &lastResult.PartialHash,
		SafeReadableFullName: safeReadableFullName,
		MediaInfo:            &lastResult.MediaInfo,
	}

	err = uc.fileStatus.Done(ctx, task.FileId, patch)
	if err != nil {
		uc.fileStatus.Failed(ctx, task.FileId, patch, uptr.String(err.Error()))
		return err
	}

	return nil
}
