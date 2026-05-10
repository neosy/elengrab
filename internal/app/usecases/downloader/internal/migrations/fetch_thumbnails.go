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

func (m *migrations) fetchThumbnails(ctx context.Context) (bool, error) {
	const (
		retryCountDefault = 3
		retryDelayDefault = 5 * time.Second
	)

	var (
		hasError bool = false
		medias   []*ddownload.File
	)

	addThumbnail := func(
		fileID uuid.UUID,
		imageData *dtypes.ImageData,
		sourceType dtypes.ThumbnailSourceType,
		mediaInfo func(thumbID uuid.UUID) *dtypes.MediaInfo,
	) error {
		var sourceURL *string
		if imageData.URL != "" {
			sourceURL = &imageData.URL
		}
		req := &dto.CreateThumbnailRequest{
			SourceType: sourceType,
			SourceURL:  sourceURL,
			ImageData:  imageData,
		}

		thubmID, err := m.usecases.thumbnail.Create(ctx, req)
		if err != nil {
			return err
		}

		mInfo := mediaInfo(thubmID)

		err = m.usecases.media.Patch(ctx, nil, fileID, &dto.FileInfoPatch{
			MediaInfo: &mInfo,
		})
		if err != nil {
			m.usecases.thumbnail.Delete(ctx, thubmID)
			return err
		}
		return nil
	}

	runWithRetry := func(run func() (*dtypes.ImageData, error)) *dtypes.ImageData {
		timer := time.NewTimer(0)
		timerStop := func() {
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
		}
		defer timerStop()

		timerStop()

		for i := range retryCountDefault {
			m.logger.Debug("Start fetch thumbnail", "attemp", i+1)
			imageData, err := run()
			if err == nil && imageData != nil {
				return imageData
			}

			timer.Reset(retryDelayDefault)

			select {
			case <-ctx.Done():
				return nil
			case <-timer.C:
			}
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

	m.logger.Debug("Found medias withow thumbnails", "count", len(medias))

	for i, media := range medias {
		m.logger.Debug("Fetching thumbnails", "index", i+1, "total", len(medias))

		select {
		case <-ctx.Done():
			return false, fmt.Errorf("context canceled: %w", ctx.Err())
		default:
		}

		mediaInfo := media.MediaInfo

		if media.MediaInfo.ThumbnailID == nil {
			imageData := runWithRetry(
				func() (*dtypes.ImageData, error) {
					return m.services.downloader.FetchThumbnail(ctx, media.MediaURL)
				},
			)
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
			filePath := m.dlStorage.Path(media.FullFileName)
			imageData, _ := m.services.ffmpeg.ExtractBestFrame(
				ctx,
				filePath,
				ffmpegsrv.WithFrameStrategy(ffmpegsrv.FrameStrategyBalanced{}),
				ffmpegsrv.WithFrameFormat(ffmpegsrv.FrameFormatWebP{}),
			)
			if imageData == nil {
				imageData = runWithRetry(
					func() (*dtypes.ImageData, error) {
						return m.services.ffmpeg.ExtractBestFrame(ctx, filePath)
					},
				)
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
