package deferred

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	ffmpegsrv "github.com/neosy/elengrab/internal/app/services/ffmpeg"
	ytdlpsrv "github.com/neosy/elengrab/internal/app/services/ytdlp"
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
			m.Logger(),
			retryCountDefault, retryDelayDefault,
			func(ctx context.Context, path string) (*dtypes.ImageData, error) {
				return m.Services().Downloader.FetchThumbnail(ctx, path, ytdlpsrv.WithRequestCookies())
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
		mutateMediaInfo func(mediaInfo *dtypes.MediaInfo, thumbID uuid.UUID),
	) error {
		thumbID, err := m.createThumbnail(ctx, imageData, sourceType)
		if err != nil {
			return err
		}

		err = m.Usecases().MediaDownload.PatchMediaInfo(
			ctx, nil, downloadID,
			func(mediaInfo *dtypes.MediaInfo) {
				mutateMediaInfo(mediaInfo, thumbID)
			},
		)
		if err != nil {
			m.Usecases().Thumbnail.Delete(ctx, thumbID)
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

	m.Logger().Debug("Found medias without thumbnails", "count", len(medias))

	var hasError bool

	for i, media := range medias {
		m.Logger().Debug("Fetching thumbnails", "index", i+1, "total", len(medias))

		select {
		case <-ctx.Done():
			return false, fmt.Errorf("context canceled: %w", ctx.Err())
		default:
		}

		if media.MediaInfo.ThumbnailID == nil {
			imageData := fetchThumbnail(ctx, media.MediaURL)
			if imageData != nil {
				sourceType := hostdetect.Detect(media.MediaURL).ThumbnailSourceType()
				if sourceType == dtypes.ThumbnailSourceTypeNone {
					sourceType = dtypes.ThumbnailSourceTypeExternal
				}

				err = createThumbnail(media.DownloadID, imageData, sourceType,
					func(mediaInfo *dtypes.MediaInfo, thumbID uuid.UUID) {
						mediaInfo.ThumbnailID = &thumbID
					},
				)
				if err != nil {
					hasError = true
				}
			}
		}

		if media.MediaInfo.FrameThumbnailID == nil && media.MediaInfo.Format.IsVideo() {
			filePath := m.DownloadsStorage().Path(media.FileFullName)
			imageData := extractBalanceThumbnail(ctx, filePath)
			if imageData == nil {
				imageData = extractThumbnail(ctx, filePath)
			}

			if imageData != nil {
				err = createThumbnail(media.DownloadID, imageData, dtypes.ThumbnailSourceTypeVideoFrame,
					func(mediaInfo *dtypes.MediaInfo, thumbID uuid.UUID) {
						mediaInfo.FrameThumbnailID = &thumbID
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
