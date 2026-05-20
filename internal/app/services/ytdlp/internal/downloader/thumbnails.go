package downloader

import (
	"context"
	"fmt"
	"time"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/executor"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/helper"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
)

func (d *Downloader) extractThumbnailFromURL(
	ctx context.Context,
	mediaURL string,
	useCookies bool,
) *dtypes.ImageData {
	imageData, err := d.FetchThumbnail(ctx, mediaURL, useCookies)
	if err != nil {
		return nil
	}

	return imageData
}

func (d *Downloader) ExtractThumbnailURL(ctx context.Context, mediaURL string, useCookies bool) (string, error) {
	if mediaURL == "" {
		return "", nil
	}

	startTime := time.Now()
	imageURL, err := d.executor.ExtractBestThumbnailURL(ctx, mediaURL, executor.WithUseCookies(useCookies))
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

func (d *Downloader) fetchYouTubeShortThumbnail(ctx context.Context, mediaURL string) (*dtypes.ImageData, error) {
	youtubeShortID, err := helper.ExtractYouTubeShortID(mediaURL)
	if err != nil {
		return nil, err
	}

	imageURL := fmt.Sprintf("https://i.ytimg.com/vi/%s/oar2.jpg", youtubeShortID)
	imageData, err := helper.FetchImage(ctx, imageURL)
	if err != nil {
		return nil, err
	}

	return imageData, nil
}

func (d *Downloader) FetchThumbnail(ctx context.Context, mediaURL string, useCookies bool) (*dtypes.ImageData, error) {
	imageData, _ := d.fetchYouTubeShortThumbnail(ctx, mediaURL)
	if imageData != nil {
		return imageData, nil
	}

	imageURL, err := d.ExtractThumbnailURL(ctx, mediaURL, useCookies)
	if err != nil {
		return nil, err
	}

	if imageURL == "" {
		return nil, nil
	}

	imageData, err = helper.FetchImage(ctx, imageURL)
	if err != nil {
		d.logger.Debug("Failed to fetch thumbnail from URL", "mediaUrl", mediaURL, "thumbnailURL", imageURL, "error", err)
		return nil, err
	}

	return imageData, nil
}
