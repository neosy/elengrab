package migrations

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	ffmpegsrv "github.com/neosy/elengrab/internal/app/services/ffmpeg"
	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (m *migrations) fillMediaInfo(ctx context.Context) (bool, error) {
	const (
		retryCountDefault = 3
		retryDelayDefault = 5 * time.Second
	)

	var (
		fetchThumbnail = makeFetchThumbnail(
			m.logger,
			retryCountDefault, retryDelayDefault,
			func(ctx context.Context, path string) (*dtypes.ImageData, error) {
				return m.services.downloader.FetchThumbnail(ctx, path, true)
			},
		)

		extractThumbnail = makeFetchThumbnail(
			m.logger,
			1, 0,
			func(ctx context.Context, path string) (*dtypes.ImageData, error) {
				return m.services.ffmpeg.ExtractBestFrame(ctx, path)
			},
		)

		extractBalanceThumbnail = makeFetchThumbnail(
			m.logger,
			1, 0,
			func(ctx context.Context, path string) (*dtypes.ImageData, error) {
				return m.services.ffmpeg.ExtractBestFrame(
					ctx,
					path,
					ffmpegsrv.WithFrameStrategy(ffmpegsrv.FrameStrategyBalanced{}),
					ffmpegsrv.WithFrameFormat(ffmpegsrv.FrameFormatWebP{}),
				)
			},
		)
	)

	medias, err := m.getAllDownloads(ctx, false,
		func(download *ddownload.MediaDownload) bool {
			if download.MediaInfo == nil {
				return true
			}
			return false
		},
	)
	if err != nil {
		return false, err
	}

	if len(medias) == 0 {
		return true, nil
	}

	m.logger.Debug("Found medias with empty mediaInfo", "count", len(medias))

	var hasError bool

	for i, media := range medias {
		m.logger.Debug("Fill media info", "index", i+1, "total", len(medias))

		select {
		case <-ctx.Done():
			return false, fmt.Errorf("context canceled: %w", ctx.Err())
		default:
		}

		var formatType dtypes.FormatType
		fileFormat := dtypes.MapFileExtToFileFormat(media.Ext)
		if fileFormat != dtypes.FileFormatNone {
			if fileFormat.IsAudio() {
				formatType = dtypes.FormatTypeAudioOnly
			} else {
				formatType = dtypes.FormatTypeVideoAudio
			}
		}

		srvMediaInfo := &dservices.MediaInfo{
			FormatType: formatType,
			Format:     fileFormat,
		}

		mediaInfoResp, _ := m.services.ffmpeg.ExtractVideoAudioInfoFromFile(
			ctx,
			m.dlStorage.Path(media.FileFullName),
			srvMediaInfo,
		)

		var (
			videoInfo *dtypes.VideoInfo
			audioInfo *dtypes.AudioInfo
		)
		if mediaInfoResp != nil {
			formatType = mediaInfoResp.FormatType
			videoInfo = mediaInfoResp.VideoInfo
			audioInfo = mediaInfoResp.AudioInfo
		}

		imageData := fetchThumbnail(ctx, media.MediaURL)

		var thumbnailID, frameThumbnailID *uuid.UUID
		if imageData != nil {
			sourceType := hostdetect.Detect(media.MediaURL).ThumbnailSourceType()
			if sourceType == dtypes.ThumbnailSourceTypeNone {
				sourceType = dtypes.ThumbnailSourceTypeExternal
			}

			thumbID, err := m.createThumbnail(ctx, imageData, sourceType)
			if err == nil {
				thumbnailID = &thumbID
			}
		}

		if fileFormat.IsVideo() {
			filePath := m.dlStorage.Path(media.FileFullName)
			imageData = extractBalanceThumbnail(ctx, filePath)
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

		err = m.usecases.download.Patch(
			ctx, nil, media.DownloadID,
			func(download *ddownload.MediaDownload) {
				download.MediaInfo = mediaInfo
			},
		)
		if err != nil {
			hasError = true
		}
	}

	if hasError {
		return false, fmt.Errorf("errors in the migration process")
	}

	return true, nil
}
