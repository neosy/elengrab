package downloader

import (
	"context"
	"fmt"
	"time"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/helper"
	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
)

func (d *Downloader) extractThumbnailFromURL(
	ctx context.Context,
	mediaURL string,
	options idto.RequestOptions,
) *dtypes.ImageData {
	imageData, err := d.FetchThumbnail(ctx, mediaURL, options)
	if err != nil {
		return nil
	}

	return imageData
}

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

func (d *Downloader) fetchYouTubeShortThumbnail(ctx context.Context, mediaURL string, options idto.RequestOptions) (*dtypes.ImageData, error) {
	youtubeShortID, err := helper.ExtractYouTubeShortID(mediaURL)
	if err != nil {
		return nil, err
	}

	imageURL := fmt.Sprintf("https://i.ytimg.com/vi/%s/oar2.jpg", youtubeShortID)
	imageData, err := helper.FetchImage(ctx, imageURL, options)
	if err != nil {
		return nil, err
	}

	return imageData, nil
}

func (d *Downloader) fetchYouTubeThumbnail(ctx context.Context, mediaURL string, options idto.RequestOptions) (*dtypes.ImageData, error) {
	youtubeID, err := helper.ExtractYouTubeID(mediaURL)
	if err != nil {
		return nil, err
	}

	// https://img.youtube.com/vi/<id>/maxresdefault.jpg		- Maximum Resolution / Max Resolution
	// https://img.youtube.com/vi/<id>/hqdefault.jpg			- High Quality
	// https://img.youtube.com/vi/<id>/mqdefault.jpg			- Medium Quality
	// https://img.youtube.com/vi/<id>/default.jpg				- Default Quality / Standard Thumbnail
	// https://img.youtube.com/vi/<id>/sddefault.jpg			- Standard Definition
	imageURL := fmt.Sprintf("https://img.youtube.com/vi/%s/hqdefault.jpg", youtubeID)
	imageData, err := helper.FetchImage(ctx, imageURL, options)
	if err != nil {
		return nil, err
	}

	return imageData, nil
}

func (d *Downloader) FetchThumbnail(
	ctx context.Context,
	mediaURL string,
	options idto.RequestOptions,
) (*dtypes.ImageData, error) {
	if hostdetect.YouTube(mediaURL) {
		imageData, _ := d.fetchYouTubeThumbnail(ctx, mediaURL, options)
		if imageData != nil {
			return imageData, nil
		}

		imageData, _ = d.fetchYouTubeShortThumbnail(ctx, mediaURL, options)
		if imageData != nil {
			return imageData, nil
		}
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
