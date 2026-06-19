package downloader

import (
	"context"
	"time"

	"github.com/google/uuid"
	ytdlpsrv "github.com/neosy/elengrab/internal/app/services/ytdlp"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

const fetchImageTimeout = 15 * time.Second

func (uc *Downloader) ExecuteRefreshMetadataTask(
	ctx context.Context,
	workerID uint64,
	task *ddownload.RefreshMetadataTask,
) error {
	fail := func() {
		err := uc.downloadStatus.Done(ctx, task.DownloadID, nil)
		if err != nil {
			uc.logger.Warn(
				"failed to update download status to done",
				"downloadID", task.DownloadID,
				"error", err,
			)
		}
		uc.logger.Error(
			"metadata refresh task failed",
			"workerID", workerID,
			"downloadID", task.DownloadID,
			"error", err,
		)
	}

	media, err := uc.download.GetByDownloadID(ctx, task.DownloadID)
	if err != nil {
		fail()
		return err
	}

	var thumbnail, frameThumbnail *dmedia.Thumbnail
	if media.MediaInfo != nil {
		if media.MediaInfo.ThumbnailID != nil {
			thumbnail, _ = uc.thumbnail.LoadByThumbID(ctx, *media.MediaInfo.ThumbnailID)
		}
		if media.MediaInfo.FrameThumbnailID != nil {
			frameThumbnail, _ = uc.thumbnail.LoadByThumbID(ctx, *media.MediaInfo.FrameThumbnailID)
		}
	}

	var (
		metadataPatch = struct {
			description *string

			mediaInfo dtypes.MediaInfo

			thumbnailData      *dtypes.ImageData
			frameThumbnailData *dtypes.ImageData
		}{}

		needPatch bool
	)

	if media.MediaInfo != nil {
		metadataPatch.mediaInfo = *media.MediaInfo.Copy()
		mediaInfo := dtypes.NewMediaInfo(media.Ext)

		if media.MediaInfo.Format != mediaInfo.Format && mediaInfo.Format != dtypes.FileFormatNone {
			needPatch = true
			metadataPatch.mediaInfo.Format = mediaInfo.Format
		}

		if media.MediaInfo.FormatType != mediaInfo.FormatType && mediaInfo.FormatType != dtypes.FormatTypeNone {
			needPatch = true
			metadataPatch.mediaInfo.FormatType = mediaInfo.FormatType
		}
	} else {
		needPatch = true
		metadataPatch.mediaInfo = *dtypes.NewMediaInfo(media.Ext)
	}

	fileMediaInfo, err := uc.ffmpegSrv.ExtractVideoAudioInfoFromFile(
		ctx,
		uc.downloadsStorage.Path(media.FileFullName),
		dservices.NewMediaInfo(media.Ext),
	)
	if err != nil {
		uc.logger.Warn(
			"Failed to extact media info from file",
			"downloadID", task.DownloadID,
			"mediaTitle", media.MediaTitle,
			"fileName", media.FileFullName,
			"error", err,
		)
		fail()
		return err
	}

	if fileMediaInfo != nil {
		if metadataPatch.mediaInfo.Format != fileMediaInfo.Format && fileMediaInfo.Format != dtypes.FileFormatNone {
			needPatch = true
			metadataPatch.mediaInfo.Format = fileMediaInfo.Format
		}

		if metadataPatch.mediaInfo.FormatType != fileMediaInfo.FormatType && fileMediaInfo.FormatType != dtypes.FormatTypeNone {
			needPatch = true
			metadataPatch.mediaInfo.FormatType = fileMediaInfo.FormatType
		}

		if metadataPatch.mediaInfo.Duration != fileMediaInfo.DurationSecondsString() && fileMediaInfo.Duration > 0 {
			needPatch = true
			metadataPatch.mediaInfo.SetDuration(fileMediaInfo.Duration)
		}

		if fileMediaInfo.VideoInfo != nil {
			if metadataPatch.mediaInfo.VideoInfo == nil {
				needPatch = true
				metadataPatch.mediaInfo.VideoInfo = fileMediaInfo.VideoInfo
			} else {
				needPatch = metadataPatch.mediaInfo.VideoInfo.Merge(*fileMediaInfo.VideoInfo) || needPatch
			}
		}

		if fileMediaInfo.AudioInfo != nil {
			if metadataPatch.mediaInfo.AudioInfo == nil {
				needPatch = true
				metadataPatch.mediaInfo.AudioInfo = fileMediaInfo.AudioInfo
			} else {
				needPatch = metadataPatch.mediaInfo.AudioInfo.Merge(*fileMediaInfo.AudioInfo) || needPatch
			}
		}
	}

	mediaInfo, err := uc.downloaderSrv.FetchInfo(ctx, media.MediaURL)
	if err != nil {
		uc.logger.Warn(
			"Failed to fetch media info",
			"downloadID", task.DownloadID,
			"mediaTitle", media.MediaTitle,
			"mediaURL", media.MediaURL,
			"error", err,
		)
	}

	if mediaInfo != nil {
		var description string
		if media.MediaDescription != nil {
			description = *media.MediaDescription
		}
		if mediaInfo.Description != "" && mediaInfo.Description != description {
			needPatch = true
			metadataPatch.description = &description
		}
	}

	thumbnailData, err := uc.downloaderSrv.FetchThumbnail(
		ctx, media.MediaURL,
		ytdlpsrv.WithRequestTimeout(fetchImageTimeout),
	)
	if err != nil {
		uc.logger.Warn(
			"failed to fetch media thumbnail",
			"downloadID", task.DownloadID,
			"mediaTitle", media.MediaTitle,
			"mediaURL", media.MediaURL,
			"error", err,
		)
	}

	if thumbnailData != nil {
		if thumbnail == nil {
			needPatch = true
			metadataPatch.thumbnailData = thumbnailData
		} else if !thumbnail.ImageDataWithSourceURL().Equal(thumbnailData) {
			needPatch = true
			metadataPatch.thumbnailData = thumbnailData
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
				"downloadID", task.DownloadID,
				"mediaTitle", media.MediaTitle,
				"fileName", media.FileFullName,
				"error", err,
			)
		}

		if frameThumbnailData != nil {
			needPatch = true
			metadataPatch.frameThumbnailData = frameThumbnailData
		}
	}

	if needPatch {
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

		patch := func(d *ddownload.MediaDownload) bool {
			if metadataPatch.description != nil {
				d.MediaDescription = metadataPatch.description
			}
			d.MediaInfo = &metadataPatch.mediaInfo
			return true
		}

		err = uc.download.Patch(ctx, nil, media.DownloadID, patch)
		if err != nil {
			fail()
			return err
		}
	}

	err = uc.downloadStatus.Done(ctx, task.DownloadID, nil)
	if err != nil {
		uc.logger.Warn(
			"failed to update download status to done",
			"downloadID", task.DownloadID,
			"error", err,
		)
		return err
	}

	uc.broadcastDownloadUpdate(ctx, task.DownloadID)

	uc.logger.Debug(
		"Metadata refresh task completed",
		"workerID", workerID,
		"downloadID", task.DownloadID,
		"title", media.MediaTitle,
	)

	return nil
}
