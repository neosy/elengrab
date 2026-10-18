package deferred

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/services/ffmpeg/ffframe"
	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

func (m *migrations) fillMediaInfo(ctx context.Context) error {
	const (
		retryCountDefault = 3
		retryDelayDefault = 5 * time.Second
	)

	var (
		fetchThumbnail = makeFetchThumbnail(
			m.Logger(),
			retryCountDefault, retryDelayDefault,
			func(ctx context.Context, path string) (*dtypes.ImageData, error) {
				return m.Services().Downloader.FetchThumbnail(ctx, path, pservices.WithDownloaderCookies())
			},
		)

		extractThumbnail = makeFetchThumbnail(
			m.Logger(),
			1, 0,
			func(ctx context.Context, path string) (*dtypes.ImageData, error) {
				return m.Services().FFMpeg.ExtractBestFrame(ctx, path)
			},
		)

		extractBalanceThumbnail = makeFetchThumbnail(
			m.Logger(),
			1, 0,
			func(ctx context.Context, path string) (*dtypes.ImageData, error) {
				return m.Services().FFMpeg.ExtractBestFrame(
					ctx,
					path,
					pservices.WithFrameStrategy(ffframe.FrameStrategyBalanced{}),
					pservices.WithFrameFormat(ffframe.FrameFormatWebP{}),
				)
			},
		)
	)

	medias, err := m.getAllDownloads(ctx, false,
		func(media *ddownload.MediaDownload) bool {
			if media.SourceType != dtypes.MediaDownloadSourceTypeURL {
				return false
			}

			if media.SourceURL() == "" {
				return false
			}

			if media.MediaInfo == nil {
				return true
			}
			return false
		},
	)
	if err != nil {
		return err
	}

	if len(medias) == 0 {
		return nil
	}

	m.Logger().Debug("Found medias with empty mediaInfo", "count", len(medias))

	var hasError bool

	for i, media := range medias {
		m.Logger().Debug("Fill media info", "index", i+1, "total", len(medias))

		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled: %w", ctx.Err())
		default:
		}

		mediaInfoResp, _ := m.Services().FFMpeg.ExtractMediaInfoWithFFprobe(
			ctx,
			m.DownloadsStorage().Path(media.FileFullName),
			pservices.NewMediaInfoFromExt(media.Ext),
		)

		var (
			fileFormat = dtypes.FileExtToFileFormat(media.Ext)
			formatType = fileFormat.FormatType()

			videoInfo *dtypes.VideoInfo
			audioInfo *dtypes.AudioInfo
		)

		if mediaInfoResp != nil {
			formatType = mediaInfoResp.FormatType
			videoInfo = mediaInfoResp.VideoInfo
			audioInfo = mediaInfoResp.AudioInfo
		}

		var thumbnailID, frameThumbnailID *uuid.UUID

		imageData := fetchThumbnail(ctx, media.SourceURL())

		if imageData != nil {
			sourceType := hostdetect.DetectPlatformType(media.SourceURL()).ThumbnailSourceType()
			if sourceType == dtypes.ThumbnailSourceTypeNone {
				sourceType = dtypes.ThumbnailSourceTypeExternal
			}

			thumbID, err := m.createThumbnail(ctx, imageData, sourceType)
			if err == nil {
				thumbnailID = &thumbID
			}
		}

		if fileFormat.IsVideo() {
			filePath := m.DownloadsStorage().Path(media.FileFullName)
			imageData := extractBalanceThumbnail(ctx, filePath)
			if imageData == nil {
				imageData = extractThumbnail(ctx, filePath)
			}

			if imageData != nil {
				thumbID, err := m.createThumbnail(ctx, imageData, dtypes.ThumbnailSourceTypeVideoFrame)
				if err == nil {
					frameThumbnailID = &thumbID
				}
			}
		}

		mediaInfo := &dtypes.MediaInfo{
			FormatType:       formatType,
			Format:           fileFormat,
			VideoInfo:        videoInfo,
			AudioInfo:        audioInfo,
			ThumbnailID:      thumbnailID,
			FrameThumbnailID: frameThumbnailID,
		}

		err = m.Usecases().MediaDownload.Patch(
			ctx, nil, media.DownloadID,
			func(download *ddownload.MediaDownload) error {
				download.MediaInfo = mediaInfo
				return nil
			},
		)
		if err != nil {
			hasError = true
		}
	}

	if hasError {
		return fmt.Errorf("errors in the migration process")
	}

	return nil
}
