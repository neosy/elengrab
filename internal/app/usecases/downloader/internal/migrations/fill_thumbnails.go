package migrations

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	ffmpegsrv "github.com/neosy/elengrab/internal/app/services/ffmpeg"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (m *migrations) fillThumbnails(ctx context.Context) (bool, error) {
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

	createThumbnail := func(
		downloadID uuid.UUID,
		imageData *dtypes.ImageData,
		sourceType dtypes.ThumbnailSourceType,
		mediaInfo func(thumbID uuid.UUID) *dtypes.MediaInfo,
	) error {
		thumbID, err := m.createThumbnail(ctx, imageData, sourceType)
		if err != nil {
			return err
		}

		mInfo := mediaInfo(thumbID)

		err = m.usecases.download.Patch(ctx, nil, downloadID, &dto.MediaDownloadInfoPatch{
			MediaInfo: &mInfo,
		})
		if err != nil {
			m.usecases.thumbnail.Delete(ctx, thumbID)
			return err
		}
		return nil
	}

	medias, err := m.getAllDownloads(ctx, false,
		func(download *ddownload.MediaDownload) bool {
			if download.MediaInfo.ThumbnailID == nil && download.MediaInfo.FrameThumbnailID == nil {
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

	m.logger.Debug("Found medias without thumbnails", "count", len(medias))

	var hasError bool

	for i, media := range medias {
		m.logger.Debug("Fetching thumbnails", "index", i+1, "total", len(medias))

		select {
		case <-ctx.Done():
			return false, fmt.Errorf("context canceled: %w", ctx.Err())
		default:
		}

		mediaInfo := media.MediaInfo

		if media.MediaInfo.ThumbnailID == nil {
			imageData := fetchThumbnail(ctx, media.MediaURL)
			if imageData != nil {
				sourceType := hostdetect.Detect(media.MediaURL).ThumbnailSourceType()
				if sourceType == dtypes.ThumbnailSourceTypeNone {
					sourceType = dtypes.ThumbnailSourceTypeExternal
				}

				err = createThumbnail(media.DownloadID, imageData, sourceType,
					func(thumbID uuid.UUID) *dtypes.MediaInfo {
						mediaInfo.ThumbnailID = &thumbID
						return mediaInfo
					},
				)
				if err != nil {
					hasError = true
				}
			}
		}

		if media.MediaInfo.FrameThumbnailID == nil && media.MediaInfo.Format.IsVideo() {
			filePath := m.dlStorage.Path(media.FileFullName)
			imageData := extractBalanceThumbnail(ctx, filePath)
			if imageData == nil {
				imageData = extractThumbnail(ctx, filePath)
			}

			if imageData != nil {
				err = createThumbnail(media.DownloadID, imageData, dtypes.ThumbnailSourceTypeVideoFrame,
					func(thumbID uuid.UUID) *dtypes.MediaInfo {
						mediaInfo.FrameThumbnailID = &thumbID
						return mediaInfo
					},
				)
				if err != nil {
					hasError = true
				}
			}
		}
	}

	if hasError {
		return false, fmt.Errorf("errors in the migration process")
	}

	return true, nil
}
