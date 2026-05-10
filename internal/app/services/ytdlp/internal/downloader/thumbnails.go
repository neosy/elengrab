package downloader

import (
	"context"
	"time"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/helper"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
)

func (d *Downloader) extractThumbnailFromURL(
	ctx context.Context,
	mediaURL string,
	onDone func(imageData *dtypes.ImageData),
) {
	imageData, err := d.FetchThumbnail(ctx, mediaURL)
	if err != nil {
		return
	}

	if imageData == nil {
		return
	}

	onDone(imageData)
}

func (d *Downloader) ExtractThumbnailURL(ctx context.Context, mediaURL string) (string, error) {
	if mediaURL == "" {
		return "", nil
	}

	startTime := time.Now()
	imageURL, err := d.executor.ExtractBestThumbnailURL(ctx, mediaURL)
	elapsed := time.Since(startTime)
	if err != nil {
		d.logger.Debug("Failed to extract thumbnail url", "mediaUrl", mediaURL, "error", err)
		return "", err
	}

	if imageURL == "" {
		d.logger.Debug("Thumbnail not found", "mediaUrl", mediaURL)
		return "", nil
	}

	d.logger.Info(
		"Thumbnail url finded",
		"mediaURL", mediaURL,
		"thumbnailUrl", imageURL,
		"elapsed", uformat.DurationFormat(elapsed),
	)

	return imageURL, nil
}

func (d *Downloader) FetchThumbnail(ctx context.Context, mediaURL string) (*dtypes.ImageData, error) {
	imageURL, err := d.ExtractThumbnailURL(ctx, mediaURL)
	if err != nil {
		return nil, err
	}

	if imageURL == "" {
		return nil, nil
	}

	imageData, err := helper.FetchImage(ctx, imageURL)
	if err != nil {
		d.logger.Debug("Failed to fetch thumbnail from URL", "mediaUrl", mediaURL, "thumbnailURL", imageURL, "error", err)
		return nil, err
	}

	return imageData, nil
}
