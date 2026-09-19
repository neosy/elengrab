package downloader

import (
	"context"
	"time"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/helper"
	"github.com/neosy/elengrab/internal/app/utils/siteimage/thumbnails"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
)

func (d *Downloader) ExtractThumbnailURL(
	ctx context.Context,
	mediaURL string,
	options idto.RequestOptions,
) (string, error) {
	if mediaURL == "" {
		return "", nil
	}

	var cookieFileName string
	if options.AllowCookies {
		cookieFileName, _ = helper.CookieFilePathFromURL(mediaURL, d.serviceOptions.CookiesDir)
	}

	startTime := time.Now()
	imageURL, err := d.executor.ExtractBestThumbnailURL(ctx, mediaURL, idto.WithUseCookies(cookieFileName))
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

func (d *Downloader) FetchThumbnail(
	ctx context.Context,
	mediaURL string,
	options idto.RequestOptions,
) (*dtypes.ImageData, error) {
	images, _ := thumbnails.FetchImages(ctx, mediaURL, options.ThumbnailFetchOptions())
	if len(images) != 0 && images[0] != nil {
		return images[0], nil
	}

	imageURL, err := d.ExtractThumbnailURL(ctx, mediaURL, options)
	if err != nil {
		return nil, err
	}

	if imageURL == "" {
		return nil, nil
	}

	imageData, err := helper.FetchImage(ctx, imageURL, options)
	if err != nil {
		d.logger.Debug("Failed to fetch thumbnail from URL", "mediaUrl", mediaURL, "thumbnailURL", imageURL, "error", err)
		return nil, err
	}

	return imageData, nil
}
