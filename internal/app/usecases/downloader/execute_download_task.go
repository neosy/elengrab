package downloader

import (
	"context"
	"fmt"
	"sync"
	"time"

	apperrors "github.com/neosy/elengrab/internal/app/errors"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

func (uc *Downloader) ExecuteDownloadTask(
	ctx context.Context,
	workerId uint64,
	task *ddownload.DownloadTask,
) error {
	if task == nil {
		uc.logger.Error("Nil pointer in function", "func", "ExecuteDownloadTask")
		return apperrors.ErrFuncParamNullPointer
	}

	err := uc.downloadStatus.Working(ctx, task.DownloadID, task.TaskID, workerId)
	if err != nil {
		uc.logger.Error("Failed update status", "error", err)
		return err
	}

	uc.broadcastDownloadUpdate(ctx, task.DownloadID)

	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
		uc.UpdateSystemInfo()
	}()

	wg.Go(func() {
		uc.fetchIcon(ctx, task.MediaUrl)
	})

	// Broadcast update download info to clients
	defer uc.broadcastDownloadUpdate(ctx, task.DownloadID)

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
			download, e := uc.download.FindByDownloadID(uc.appCtx, nil, task.DownloadID)
			if e == nil && download != nil {
				uc.downloadStatus.Failed(uc.appCtx, task.DownloadID, nil, uptr.String(err.Error()))
			}
			uc.dlStateCache.Delete(uc.appCtx, task.DownloadID)
			return ctx.Err()
		}

		uc.logger.Error(
			"Failed to download",
			"error", err,
		)

		uc.downloadStatus.Failed(ctx, task.DownloadID, nil, uptr.String(err.Error()))

		return err
	}

	var (
		lastResult, resultBeforeBroadcast, resultProgressBeforeBroadcast *dservices.DownloaderResult

		channelProcess, thumbnailProcess sync.Once

		mediaInfoStorage dtypes.MediaInfo
	)

	mediaInfo := func(srvMediaInfo *dservices.MediaInfo) *dtypes.MediaInfo {
		if srvMediaInfo == nil {
			return nil
		}

		mediaInfo := srvMediaInfo.MediaInfoDomainPtr()

		mediaInfo.ThumbnailID = mediaInfoStorage.ThumbnailID
		mediaInfo.FrameThumbnailID = mediaInfoStorage.FrameThumbnailID

		return mediaInfo
	}

	for r := range resultCh {
		if r == nil {
			uc.logger.Error("Received nil result from result channel")
			uc.downloadStatus.Failed(uc.appCtx, task.DownloadID, nil, new(apperrors.ErrDownloaderEmptyResponse.Error()))
			uc.dlStateCache.Delete(uc.appCtx, task.DownloadID)
			return apperrors.ErrDownloaderEmptyResponse
		}

		if r.Error != nil {
			// The context was canceled
			if ctx.Err() != nil {
				download, e := uc.download.FindByDownloadID(uc.appCtx, nil, task.DownloadID)
				if e == nil && download != nil {
					uc.downloadStatus.Failed(uc.appCtx, task.DownloadID, nil, new(r.Error.Error()))
				}
				uc.dlStateCache.Delete(uc.appCtx, task.DownloadID)
				return ctx.Err()
			}

			uc.logger.Error(
				"Failed to download",
				"downloadID", task.DownloadID,
				"error", r.Error,
			)

			patch := func(download *ddownload.MediaDownload) {
				if download == nil || lastResult == nil {
					return
				}

				download.ChannelID = lastResult.ChannelID
				if lastResult.MediaTitle != "" {
					download.MediaTitle = lastResult.MediaTitle
				}
				if lastResult.FileExt != "" {
					download.Ext = lastResult.FileExt
				}
				if lastResult.Filesize != nil && *lastResult.Filesize != 0 {
					download.FileSize = lastResult.Filesize
				}
			}
			uc.downloadStatus.Failed(ctx, task.DownloadID, patch, new(r.Error.Error()))
			return r.Error
		}

		state, err := uc.dlStateCache.FindByDownloadID(ctx, nil, task.DownloadID)
		if err != nil {
			uc.logger.Error(
				"Failed to download",
				"action", "Find by downloadID",
				"error", err,
			)
			uc.downloadStatus.Failed(ctx, task.DownloadID, nil, new(err.Error()))
			return err
		}

		lastResult = r

		// Adding a record to the YouTube Channel table
		if lastResult.ChannelID != nil && lastResult.Channel != nil {
			channelProcess.Do(func() {
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
			})
		}

		if lastResult.Thumbnail != nil {
			thumbnailProcess.Do(func() {
				req := &dto.CreateThumbnailRequest{
					SourceType: dtypes.ThumbnailSourceTypeExternal,
					SourceURL:  &lastResult.Thumbnail.URL,
					ImageData:  lastResult.Thumbnail,
				}
				id, err := uc.thumbnail.Create(ctx, req)
				if err == nil {
					mediaInfoStorage.ThumbnailID = &id
				}
			})
		}

		state.InitFromDownloaderResult(lastResult, mediaInfo(lastResult.MediaInfo))
		uc.dlStateCache.Save(
			ctx,
			state,
		)

		if resultProgressBeforeBroadcast == nil {
			resultProgressBeforeBroadcast = lastResult
		}

		if lastResult.MetadataChanged(resultBeforeBroadcast) {
			uc.broadcastDownloadUpdate(ctx, task.DownloadID)
			resultBeforeBroadcast = lastResult
		} else if lastResult.ProgressChanged(resultProgressBeforeBroadcast) {
			uc.broadcastDownloadProgressUpdate(ctx, task.DownloadID)
			resultProgressBeforeBroadcast = lastResult
		}
	}

	if mediaInfo != nil && lastResult.ThumbnailVideoFrame != nil {
		req := &dto.CreateThumbnailRequest{
			SourceType: dtypes.ThumbnailSourceTypeVideoFrame,
			ImageData:  lastResult.ThumbnailVideoFrame,
		}
		id, err := uc.thumbnail.Create(ctx, req)
		if err == nil {
			mediaInfoStorage.FrameThumbnailID = &id
		}
	}

	patch := func(download *ddownload.MediaDownload) {
		download.ChannelID = lastResult.ChannelID
		download.MediaTitle = lastResult.MediaTitle
		download.MediaDescription = lastResult.MediaDescription
		download.FileName = lastResult.Filename
		download.Ext = lastResult.FileExt
		download.FileFullName = lastResult.FileFullName
		download.FileSize = lastResult.Filesize
		download.PartialHash = lastResult.PartialHash
		download.SafeReadableFullName = fmt.Sprintf("%s.%s", nfasthttp.SanitizeFileName(lastResult.MediaTitle), lastResult.FileExt)
		download.MediaInfo = mediaInfo(lastResult.MediaInfo)
	}

	err = uc.downloadStatus.Done(ctx, task.DownloadID, patch)
	if err != nil {
		uc.downloadStatus.Failed(ctx, task.DownloadID, patch, new(err.Error()))
		return err
	}

	return nil
}
