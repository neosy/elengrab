package dlexecutor

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	apperrors "github.com/neosy/elengrab/internal/app/errors"
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_executor/types"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
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
		return resultCh, nil
	}

	// The context was canceled
	if ctx.Err() != nil {
		uc.logger.Debug(
			"Failed to download: The context was canceled",
			"error", err,
		)
		return nil, ctx.Err()
	}

	uc.logger.Error(
		"Failed to download",
		"error", err,
	)

	return nil, err
}

func (uc *Executor) processDownloadResults(
	ctx context.Context,
	task *ddownload.DownloadTask,
	resultCh <-chan *dservices.DownloaderResult,
) (*types.ProcessedDownload, error) {
	var (
		lastResult, resultBeforeBroadcast, resultProgressBeforeBroadcast *dservices.DownloaderResult

		channelProcess, thumbnailProcess sync.Once
		thumbnailIDs                     types.ThumbnailIDs
		channelID                        uuid.UUID
	)

	for r := range resultCh {
		if r == nil {
			uc.logger.Error("Received nil result from result channel")
			return nil, apperrors.ErrDownloaderEmptyResponse
		}

		if r.Error != nil {
			uc.logger.Error(
				"Failed to download",
				"downloadID", task.DownloadID,
				"error", r.Error,
			)
			if lastResult != nil {
				state, _ := uc.download.FindState(ctx, task.DownloadID)
				return uc.mappers.MapDownloadResultToProcessedDownload(lastResult, state, channelID, thumbnailIDs), r.Error
			}
			return nil, r.Error
		}

		lastResult = r

		state, err := uc.download.GetOrCreateState(ctx, task.DownloadID)
		if err != nil {
			return nil, err
		}

		if state.Download != nil && state.Download.MediaInfo != nil {
			thumbnailIDs.ThumbnailID = state.Download.MediaInfo.ThumbnailID
			thumbnailIDs.FrameThumbnailID = state.Download.MediaInfo.FrameThumbnailID
		}

		// Adding a record to the YouTube Channel table
		if lastResult.Channel != nil {
			channelProcess.Do(func() {
				channelSource := lastResult.Channel
				channel, _ := uc.channel.FindByExternalChannelIDNoCache(ctx, channelSource.ChannelID, channelSource.Platform)
				if channel != nil {
					uc.channel.Patch(ctx, channel.ChannelID,
						func(c *dmedia.Channel) error {
							if c.EqualSource(channelSource) {
								return nil
							}

							c.UpdateFromSource(channelSource)

							return nil
						},
					)
				} else {
					newChannel := dmedia.NewChannelFromSource(lastResult.Channel)

					err := uc.channel.Create(ctx, newChannel)
					if err == nil {
						channel, _ = uc.channel.FindByChannelID(ctx, newChannel.ChannelID)
					}
				}

				if channel != nil {
					channelID = channel.ChannelID
				}
			})
		}

		if lastResult.Thumbnail != nil && thumbnailIDs.ThumbnailID == nil {
			thumbnailProcess.Do(func() {
				req := &dto.CreateThumbnailRequest{
					SourceType: dtypes.ThumbnailSourceTypeExternal,
					SourceURL:  &lastResult.Thumbnail.URL,
					ImageData:  lastResult.Thumbnail,
				}
				id, err := uc.thumbnail.Create(ctx, req)
				if err == nil {
					thumbnailIDs.ThumbnailID = &id
				}
			})
		}

		uc.download.PatchState(
			ctx, task.DownloadID,
			func(state *ddownload.DownloadState) error {
				if lastResult != nil {
					var mediaInfo *dtypes.MediaInfo
					if lastResult.MediaInfo != nil {
						mediaInfo = uc.mappers.MapMediaInfoDomain(lastResult.MediaInfo, thumbnailIDs)
					}
					uc.mappers.MapDownloaderResultToState(state, lastResult, mediaInfo, channelID)
				}
				return nil
			},
		)

		if resultProgressBeforeBroadcast == nil {
			resultProgressBeforeBroadcast = lastResult
		}

		if lastResult.MetadataChanged(resultBeforeBroadcast) {
			uc.broadcaster.DownloadUpdate(ctx, task.DownloadID)
			resultBeforeBroadcast = lastResult
		} else if lastResult.ProgressChanged(resultProgressBeforeBroadcast) {
			uc.broadcaster.DownloadProgressUpdate(ctx, task.DownloadID)
			resultProgressBeforeBroadcast = lastResult
		}
	}

	if lastResult == nil {
		err := fmt.Errorf("failed to download: no result returned")
		uc.logger.Error(
			"Failed to download",
			"error", err,
		)
		return nil, err
	}

	if lastResult.ThumbnailVideoFrame != nil && thumbnailIDs.FrameThumbnailID == nil {
		req := &dto.CreateThumbnailRequest{
			SourceType: dtypes.ThumbnailSourceTypeVideoFrame,
			ImageData:  lastResult.ThumbnailVideoFrame,
		}
		id, err := uc.thumbnail.Create(ctx, req)
		if err == nil {
			thumbnailIDs.FrameThumbnailID = &id
		}
	}

	state, _ := uc.download.FindState(ctx, task.DownloadID)

	return uc.mappers.MapDownloadResultToProcessedDownload(lastResult, state, channelID, thumbnailIDs), nil
}
