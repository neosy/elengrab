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
	if mediaURL == "" {
		return
	}

	startTime := time.Now()
	thumbnailURL, err := d.executor.ExtractBestThumbnailURL(ctx, mediaURL)
	elapsed := time.Since(startTime)
	if err != nil {
		d.logger.Debug("Failed to fetch thumbnail url", "mediaUrl", mediaURL, "error", err)
		return
	}

	if thumbnailURL == "" {
		d.logger.Debug("Thumbnail not found", "mediaUrl", mediaURL)
		return
	}

	d.logger.Info(
		"Thumbnail url finded",
		"mediaURL", mediaURL,
		"thumbnailUrl", thumbnailURL,
		"elapsed", uformat.DurationFormat(elapsed),
	)

	imageData, err := helper.FetchImage(ctx, thumbnailURL)
	if err != nil {
		d.logger.Debug("Failed to fetch thumbnail from URL", "mediaUrl", mediaURL, "thumbnailURL", thumbnailURL, "error", err)
		return
	}

	onDone(imageData)
}
