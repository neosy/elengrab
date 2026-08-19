package dlexecutor

import (
	"context"
	"fmt"
	"sync"
	"time"

	apperrors "github.com/neosy/elengrab/internal/app/errors"
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_executor/types"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

func (uc *Executor) startDownload(
	ctx context.Context,
	task *ddownload.DownloadTask,
) (<-chan *dservices.DownloaderResult, error) {
	resultCh, err := uc.downloaderSrv.Download(
		ctx,
		task.MediaUrl,
		uc.mappers.MapDownloadOptionsDomainToService(task.Options),
	)

	if err == nil {
		return resultCh, err
	}

	// The context was canceled
	if ctx.Err() != nil {
		uc.logger.Debug(
			"Failed to download: The context was canceled",
			"error", err,
		)
		download, e := uc.download.FindByDownloadIDNoCache(uc.appCtx, task.DownloadID)
		if e == nil && download != nil {
			uc.downloadStatus.Failed(uc.appCtx, task.DownloadID, nil, uptr.String(err.Error()))
		}
		uc.dlStateCache.Delete(uc.appCtx, task.DownloadID)
		return nil, ctx.Err()
	}

	uc.logger.Error(
		"Failed to download",
		"error", err,
	)

	uc.downloadStatus.Failed(ctx, task.DownloadID, nil, uptr.String(err.Error()))

	return nil, err
}

func (uc *Executor) processDownloadResults(
	ctx context.Context,
	task *ddownload.DownloadTask,
	resultCh <-chan *dservices.DownloaderResult,
) (*types.ProcessedDownload, error) {
	var (
		result, resultBeforeBroadcast, resultProgressBeforeBroadcast *dservices.DownloaderResult

		channelProcess, thumbnailProcess sync.Once
		thumbnailsIDs                    types.ThumbnailIDs
	)

	for r := range resultCh {
		if r == nil {
			uc.logger.Error("Received nil result from result channel")
			uc.downloadStatus.Failed(uc.appCtx, task.DownloadID, nil, new(apperrors.ErrDownloaderEmptyResponse.Error()))
			uc.dlStateCache.Delete(uc.appCtx, task.DownloadID)
			return nil, apperrors.ErrDownloaderEmptyResponse
		}

		if r.Error != nil {
			// The context was canceled
			if ctx.Err() != nil {
				download, e := uc.download.FindByDownloadIDNoCache(uc.appCtx, task.DownloadID)
				if e == nil && download != nil {
					uc.downloadStatus.Failed(uc.appCtx, task.DownloadID, nil, new(r.Error.Error()))
				}
				uc.dlStateCache.Delete(uc.appCtx, task.DownloadID)
				return nil, ctx.Err()
			}

			uc.logger.Error(
				"Failed to download",
				"downloadID", task.DownloadID,
				"error", r.Error,
			)

			patch := func(download *ddownload.MediaDownload) {
				if download == nil || result == nil {
					return
				}

				download.ChannelID = result.ChannelID
				if result.MediaTitle != "" {
					download.MediaTitle = result.MediaTitle
				}
				if result.FileExt != "" {
					download.Ext = result.FileExt
				}
				if result.Filesize != nil && *result.Filesize != 0 {
					download.FileSize = result.Filesize
				}
			}
			uc.downloadStatus.Failed(ctx, task.DownloadID, patch, new(r.Error.Error()))
			return nil, r.Error
		}

		state, err := uc.dlStateCache.FindByDownloadID(ctx, task.DownloadID)
		if err != nil {
			uc.logger.Error(
				"Failed to download",
				"action", "Find by downloadID",
				"error", err,
			)
			uc.downloadStatus.Failed(ctx, task.DownloadID, nil, new(err.Error()))
			return nil, err
		}

		result = r

		// Adding a record to the YouTube Channel table
		if result.ChannelID != nil && result.Channel != nil {
			channelProcess.Do(func() {
				channel, _ := uc.ytChannel.FindByChannelID(ctx, *result.ChannelID)
				if channel != nil {
					if time.Since(channel.UpdatedAt) > uc.channelUpdateInterval {
						channel.InitFromChannel(result.Channel)
						uc.ytChannel.Update(ctx, channel)
					}
				} else {
					channel := &dmedia.YoutubeChannel{
						ChannelID: *result.ChannelID,
					}
					channel.InitFromChannel(result.Channel)
					uc.ytChannel.Create(ctx, channel)
				}
			})
		}

		if result.Thumbnail != nil {
			thumbnailProcess.Do(func() {
				req := &dto.CreateThumbnailRequest{
					SourceType: dtypes.ThumbnailSourceTypeExternal,
					SourceURL:  &result.Thumbnail.URL,
					ImageData:  result.Thumbnail,
				}
				id, err := uc.thumbnail.Create(ctx, req)
				if err == nil {
					thumbnailsIDs.ThumbnailID = &id
				}
			})
		}

		mediaInfo := uc.mappers.MapMediaInfoDomain(result.MediaInfo, thumbnailsIDs)
		state.InitFromDownloaderResult(result, mediaInfo)
		uc.dlStateCache.Save(
			ctx,
			state,
		)

		if resultProgressBeforeBroadcast == nil {
			resultProgressBeforeBroadcast = result
		}

		if result.MetadataChanged(resultBeforeBroadcast) {
			uc.broadcaster.DownloadUpdate(ctx, task.DownloadID)
			resultBeforeBroadcast = result
		} else if result.ProgressChanged(resultProgressBeforeBroadcast) {
			uc.broadcaster.DownloadProgressUpdate(ctx, task.DownloadID)
			resultProgressBeforeBroadcast = result
		}
	}

	if result == nil {
		err := fmt.Errorf("failed to download: no result returned")
		uc.logger.Error(
			"Failed to download",
			"error", err,
		)
		uc.downloadStatus.Failed(ctx, task.DownloadID, nil, new(err.Error()))
		return nil, err
	}

	if result.ThumbnailVideoFrame != nil {
		req := &dto.CreateThumbnailRequest{
			SourceType: dtypes.ThumbnailSourceTypeVideoFrame,
			ImageData:  result.ThumbnailVideoFrame,
		}
		id, err := uc.thumbnail.Create(ctx, req)
		if err == nil {
			thumbnailsIDs.FrameThumbnailID = &id
		}
	}

	return &types.ProcessedDownload{
		Result:       result,
		ThumbnailIDs: thumbnailsIDs,
	}, nil
}
