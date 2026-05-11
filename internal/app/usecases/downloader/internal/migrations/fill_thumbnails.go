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
		hasError bool = false
		medias   []*ddownload.File
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

	addThumbnail := func(
		fileID uuid.UUID,
		imageData *dtypes.ImageData,
		sourceType dtypes.ThumbnailSourceType,
		mediaInfo func(thumbID uuid.UUID) *dtypes.MediaInfo,
	) error {
		thumbID, err := m.createThumbnail(ctx, imageData, sourceType)
		if err != nil {
			return err
		}

		mInfo := mediaInfo(thumbID)

		err = m.usecases.media.Patch(ctx, nil, fileID, &dto.FileInfoPatch{
			MediaInfo: &mInfo,
		})
		if err != nil {
			m.usecases.thumbnail.Delete(ctx, thumbID)
			return err
		}
		return nil
	}

	err := m.usecases.media.GetAll(ctx, false,
		func(media *ddownload.File) error {
			if media == nil || media.MediaInfo == nil {
				return nil
			}

			if media.MediaURL == "" {
				return nil
			}

			if media.MediaInfo.ThumbnailID != nil || media.MediaInfo.FrameThumbnailID != nil {
				return nil
			}

			medias = append(medias, media)

			return nil
		},
	)
	if err != nil {
		return false, err
	}

	if len(medias) == 0 {
		return true, nil
	}

	m.logger.Debug("Found medias without thumbnails", "count", len(medias))

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

				err = addThumbnail(media.FileID, imageData, sourceType,
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
				err = addThumbnail(media.FileID, imageData, dtypes.ThumbnailSourceTypeVideoFrame,
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
