package downloader

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	"github.com/neosy/elengrab/internal/pkg/nfasthttp"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

func (uc *YouTubeDownloader) ExecuteDownloadTask(
	ctx context.Context,
	workerId uint64,
	task *ddownload.DownloadTask,
) error {
	if task == nil {
		uc.logger.Error("Nil pointer in function", "func", "ExecuteDownloadTask")
		return errors.New("function parameter is a null pointer")
	}

	err := uc.fileStatus.Working(ctx, task.FileID, task.TaskID, workerId)
	if err != nil {
		uc.logger.Error("Failed update status", "error", err)
		return err
	}

	uc.broadcastFileUpdate(ctx, task.FileID)

	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
		uc.UpdateSystemInfo()
	}()

	wg.Go(func() {
		uc.fetchIcon(ctx, task.MediaUrl)
	})

	// Broadcast update file info to clients
	defer uc.broadcastFileUpdate(ctx, task.FileID)

	resultCh, err := uc.downloaderSrv.Download(
		ctx,
		task.MediaUrl,
		uc.mappers.MapDownloadOptionsDomainToService(task.Options),
	)
	if err != nil {
		// The context was canceled
		if ctx.Err() != nil {
			uc.logger.Debug(
				"Failed to download: The context was canceled",
				"error", err,
			)
			file, e := uc.file.FindByFileID(uc.appCtx, nil, task.FileID)
			if e == nil && file != nil {
				uc.fileStatus.Failed(uc.appCtx, task.FileID, nil, uptr.String(err.Error()))
			}
			uc.dlStateCache.Delete(uc.appCtx, task.FileID)
			return ctx.Err()
		}

		uc.logger.Error(
			"Failed to download",
			"error", err,
		)

		uc.fileStatus.Failed(ctx, task.FileID, nil, uptr.String(err.Error()))

		return err
	}

	var lastResult, resultBeforeBroadcast, resultProgressBeforeBroadcast *ddownload.DownloadResult
	for r := range resultCh {
		if r.Error != nil {
			// The context was canceled
			if ctx.Err() != nil {
				file, e := uc.file.FindByFileID(uc.appCtx, nil, task.FileID)
				if e == nil && file != nil {
					uc.fileStatus.Failed(uc.appCtx, task.FileID, nil, uptr.String(r.Error.Error()))
				}
				uc.dlStateCache.Delete(uc.appCtx, task.FileID)
				return ctx.Err()
			}

			uc.logger.Error(
				"Failed to download",
				"fileId", task.FileID,
				"error", r.Error,
			)

			var patch *dto.FileInfoPatch
			if lastResult != nil {
				patch = &dto.FileInfoPatch{
					YoutubeChannelID: &lastResult.ChannelID,
				}
				if lastResult.MediaTitle != "" {
					patch.MediaTitle = &lastResult.MediaTitle
				}
				if lastResult.FileExt != "" {
					patch.Ext = &lastResult.FileExt
				}
				if lastResult.Filesize != nil && *lastResult.Filesize != 0 {
					patch.FileSize = &lastResult.Filesize
				}
			}
			uc.fileStatus.Failed(ctx, task.FileID, patch, uptr.String(r.Error.Error()))
			return r.Error
		}

		state, err := uc.dlStateCache.FindByFileID(ctx, nil, task.FileID)
		if err != nil {
			uc.logger.Error(
				"Failed to download",
				"action", "Find by fileId",
				"error", err,
			)
			uc.fileStatus.Failed(ctx, task.FileID, nil, uptr.String(err.Error()))
			return err
		}

		lastResult = r

		// Adding a record to the YouTube Channel table
		if lastResult != nil && lastResult.ChannelID != nil && lastResult.Channel != nil {
			channel, _ := uc.ytChannel.FindByChannelID(ctx, *lastResult.ChannelID)
			if channel != nil {
				if time.Since(channel.UpdatedAt) > uc.channelUpdateInterval {
					channel.InitFromChannel(lastResult.Channel)
					uc.ytChannel.Update(ctx, channel)
				}
			} else {
				channel := &dmedia.YoutubeChannel{
					ChannelID: *lastResult.ChannelID,
				}
				channel.InitFromChannel(lastResult.Channel)
				uc.ytChannel.Create(ctx, channel)
			}
		}

		state.InitFromDownloadResult(r)
		uc.dlStateCache.Save(
			ctx,
			state,
		)

		if resultProgressBeforeBroadcast == nil {
			resultProgressBeforeBroadcast = r
		}

		if r.MetadataChanged(resultBeforeBroadcast) {
			uc.broadcastFileUpdate(ctx, task.FileID)
			resultBeforeBroadcast = r
		} else if r.ProgressChanged(resultProgressBeforeBroadcast) {
			uc.broadcastFileProgressUpdate(ctx, task.FileID)
			resultProgressBeforeBroadcast = r
		}
	}

	if lastResult == nil {
		return fmt.Errorf("service downloader returned an empty value")
	}

	safeReadableFullName := uptr.String(
		fmt.Sprintf("%s.%s", nfasthttp.SanitizeFileName(lastResult.MediaTitle), lastResult.FileExt),
	)

	patch := &dto.FileInfoPatch{
		YoutubeChannelID:     &lastResult.ChannelID,
		MediaTitle:           &lastResult.MediaTitle,
		FileName:             &lastResult.Filename,
		Ext:                  &lastResult.FileExt,
		FullName:             &lastResult.FileFullName,
		FileSize:             &lastResult.Filesize,
		PartialHash:          &lastResult.PartialHash,
		SafeReadableFullName: safeReadableFullName,
		MediaInfo:            &lastResult.MediaInfo,
	}

	err = uc.fileStatus.Done(ctx, task.FileID, patch)
	if err != nil {
		uc.fileStatus.Failed(ctx, task.FileID, patch, uptr.String(err.Error()))
		return err
	}

	return nil
}
