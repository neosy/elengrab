package dlexecutor

import (
	"context"

	"github.com/google/uuid"
	apperrors "github.com/neosy/elengrab/internal/app/errors"
	ytdlpsrv "github.com/neosy/elengrab/internal/app/services/ytdlp"
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/helper"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type refreshMetadataPatch struct {
	title       *string
	description *string

	mediaInfo dtypes.MediaInfo

	thumbnailData      *dtypes.ImageData
	frameThumbnailData *dtypes.ImageData

	channelID uuid.UUID

	needPatch bool
}

func (uc *Executor) collectMetadata(
	ctx context.Context,
	media *ddownload.MediaDownload,
) (*refreshMetadataPatch, error) {
	var thumbnail, frameThumbnail *dmedia.Thumbnail
	if media.MediaInfo != nil {
		if media.MediaInfo.ThumbnailID != nil {
			thumbnail, _ = uc.thumbnail.LoadByThumbID(ctx, *media.MediaInfo.ThumbnailID)
		}
		if media.MediaInfo.FrameThumbnailID != nil {
			frameThumbnail, _ = uc.thumbnail.LoadByThumbID(ctx, *media.MediaInfo.FrameThumbnailID)
		}
	}

	patch := &refreshMetadataPatch{}

	if media.MediaInfo != nil {
		patch.mediaInfo = *media.MediaInfo.Clone()
		mediaInfo := dtypes.NewMediaInfo(media.Ext)

		if media.MediaInfo.Format != mediaInfo.Format && mediaInfo.Format != dtypes.FileFormatNone {
			patch.needPatch = true
			patch.mediaInfo.Format = mediaInfo.Format
		}

		if media.MediaInfo.FormatType != mediaInfo.FormatType && mediaInfo.FormatType != dtypes.FormatTypeNone {
			patch.needPatch = true
			patch.mediaInfo.FormatType = mediaInfo.FormatType
		}
	} else {
		patch.needPatch = true
		patch.mediaInfo = *dtypes.NewMediaInfo(media.Ext)
	}

	fileMediaInfo, err := uc.ffmpegSrv.ExtractVideoAudioInfoFromFile(
		ctx,
		uc.downloadsStorage.Path(media.FileFullName),
		dservices.NewMediaInfo(media.Ext),
	)
	if err != nil {
		uc.logger.Warn(
			"Failed to extact media info from file",
			"downloadID", media.DownloadID,
			"mediaTitle", media.MediaTitle,
			"fileName", media.FileFullName,
			"error", err,
		)
		return nil, err
	}

	if fileMediaInfo != nil {
		if patch.mediaInfo.Format != fileMediaInfo.Format && fileMediaInfo.Format != dtypes.FileFormatNone {
			patch.needPatch = true
			patch.mediaInfo.Format = fileMediaInfo.Format
		}

		if patch.mediaInfo.FormatType != fileMediaInfo.FormatType && fileMediaInfo.FormatType != dtypes.FormatTypeNone {
			patch.needPatch = true
			patch.mediaInfo.FormatType = fileMediaInfo.FormatType
		}

		if patch.mediaInfo.DurationText != fileMediaInfo.DurationSecondsString() && fileMediaInfo.Duration > 0 {
			patch.needPatch = true
			patch.mediaInfo.SetDuration(fileMediaInfo.Duration)
		}

		if fileMediaInfo.VideoInfo != nil {
			if patch.mediaInfo.VideoInfo == nil {
				patch.needPatch = true
				patch.mediaInfo.VideoInfo = fileMediaInfo.VideoInfo
			} else {
				patch.needPatch = patch.mediaInfo.VideoInfo.Merge(*fileMediaInfo.VideoInfo) || patch.needPatch
			}
		}

		if fileMediaInfo.AudioInfo != nil {
			if patch.mediaInfo.AudioInfo == nil {
				patch.needPatch = true
				patch.mediaInfo.AudioInfo = fileMediaInfo.AudioInfo
			} else {
				patch.needPatch = patch.mediaInfo.AudioInfo.Merge(*fileMediaInfo.AudioInfo) || patch.needPatch
			}
		}
	}

	mediaInfo, err := uc.downloaderSrv.FetchInfo(ctx, media.MediaURL)
	if err != nil {
		uc.logger.Warn(
			"Failed to fetch media info",
			"downloadID", media.DownloadID,
			"mediaTitle", media.MediaTitle,
			"mediaURL", media.MediaURL,
			"error", err,
		)
	}

	if mediaInfo != nil {
		if mediaInfo.Title != media.MediaTitle {
			patch.needPatch = true
			patch.title = &mediaInfo.Title
		}

		if helper.ValuesEqual(media.MediaDescription, &mediaInfo.Description) {
			patch.needPatch = true
			patch.description = &mediaInfo.Description
		}
	}

	channelSource, err := uc.downloaderSrv.FetchChannelInfoWithCookieFallback(ctx, media.MediaURL)
	if err != nil {
		uc.logger.Warn(
			"Failed to fetch channel info",
			"downloadID", media.DownloadID,
			"mediaTitle", media.MediaTitle,
			"mediaURL", media.MediaURL,
			"error", err,
		)
	}

	if channelSource != nil && channelSource.ChannelID != "" && channelSource.Platform != "" {
		var channelID uuid.UUID
		channel, _ := uc.channel.FindByExternalChannelIDNoCache(ctx, channelSource.ChannelID, channelSource.Platform)
		if channel != nil {
			channelID = channelID

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
			newChannel := dmedia.NewChannelFromSource(channelSource)
			err := uc.channel.Create(ctx, newChannel)
			if err == nil {
				channelID = newChannel.ChannelID
			}
		}

		if channelID != uuid.Nil && media.ChannelID == uuid.Nil {
			patch.needPatch = true
			patch.channelID = channelID
		}
	}

	thumbnailData, err := uc.downloaderSrv.FetchThumbnail(
		ctx, media.MediaURL,
		ytdlpsrv.WithRequestTimeout(fetchImageTimeout),
	)
	if err != nil {
		uc.logger.Warn(
			"failed to fetch media thumbnail",
			"downloadID", media.DownloadID,
			"mediaTitle", media.MediaTitle,
			"mediaURL", media.MediaURL,
			"error", err,
		)
	}

	if thumbnailData != nil {
		if thumbnail == nil {
			patch.needPatch = true
			patch.thumbnailData = thumbnailData
		} else if !thumbnail.ImageDataWithSourceURL().Equal(thumbnailData) {
			patch.needPatch = true
			patch.thumbnailData = thumbnailData
		}
	}

	if frameThumbnail == nil {
		frameThumbnailData, err := uc.ffmpegSrv.ExtractBestFrame(
			ctx,
			uc.downloadsStorage.Path(media.FileFullName),
		)
		if err != nil {
			uc.logger.Warn(
				"failed to extract frame from file",
				"downloadID", media.DownloadID,
				"mediaTitle", media.MediaTitle,
				"fileName", media.FileFullName,
				"error", err,
			)
		}

		if frameThumbnailData != nil {
			patch.needPatch = true
			patch.frameThumbnailData = frameThumbnailData
		}
	}

	return patch, nil
}

func (uc *Executor) applyMetadataPatch(
	ctx context.Context,
	downloadID uuid.UUID,
	metadataPatch *refreshMetadataPatch,
) error {
	if downloadID != uuid.Nil && metadataPatch == nil {
		return apperrors.ErrFuncParamNullPointer
	}

	if !metadataPatch.needPatch {
		return nil
	}

	if metadataPatch.thumbnailData != nil {
		req := &dto.CreateThumbnailRequest{
			Variant:    dtypes.ThumbnailVariantOriginal,
			SourceType: dtypes.ThumbnailSourceTypeExternal,
			SourceURL:  &metadataPatch.thumbnailData.URL,
			ImageData:  metadataPatch.thumbnailData,
		}
		thumbnailID, err := uc.thumbnail.Create(ctx, req)
		if err != nil {
			uc.logger.Warn(
				"failed to create thumbnail",
				"sourceType", dtypes.ThumbnailSourceTypeExternal,
				"sourceURL", metadataPatch.thumbnailData.URL,
				"error", err,
			)
		}
		if thumbnailID != uuid.Nil {
			metadataPatch.mediaInfo.ThumbnailID = &thumbnailID
		}
	}

	if metadataPatch.frameThumbnailData != nil {
		req := &dto.CreateThumbnailRequest{
			Variant:    dtypes.ThumbnailVariantOriginal,
			SourceType: dtypes.ThumbnailSourceTypeVideoFrame,
			ImageData:  metadataPatch.frameThumbnailData,
		}
		thumbnailID, err := uc.thumbnail.Create(ctx, req)
		if err != nil {
			uc.logger.Warn(
				"failed to create thumbnail",
				"sourceType", dtypes.ThumbnailSourceTypeVideoFrame,
				"sourceURL", metadataPatch.thumbnailData.URL,
				"error", err,
			)
		}
		if thumbnailID != uuid.Nil {
			metadataPatch.mediaInfo.FrameThumbnailID = &thumbnailID
		}
	}

	patch := func(d *ddownload.MediaDownload) error {
		if metadataPatch.title != nil {
			if d.MediaTitle == d.MediaTitleOriginal {
				d.MediaTitle = *metadataPatch.title
			}
			d.MediaTitleOriginal = *metadataPatch.title
		}

		if metadataPatch.description != nil {
			if helper.ValuesEqual(d.MediaDescription, d.MediaDescriptionOriginal) {
				d.MediaDescription = metadataPatch.description
			}
			d.MediaDescriptionOriginal = metadataPatch.description
		}

		if metadataPatch.channelID != uuid.Nil {
			d.ChannelID = metadataPatch.channelID
		}

		d.MediaInfo = &metadataPatch.mediaInfo

		return nil
	}

	return uc.download.Patch(ctx, nil, downloadID, patch)
}
