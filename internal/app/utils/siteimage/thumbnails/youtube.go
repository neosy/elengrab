package thumbnails

import (
	"context"
	"fmt"
	"slices"

	"github.com/neosy/elengrab/internal/app/utils/siteimage/internal/fetch"
	"github.com/neosy/elengrab/internal/app/utils/websource/youtube"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

var (
	youTubeShortThumbnailURLTemplateList = []string{
		"https://i.ytimg.com/vi/%s/oardefault.jpg",
		"https://i.ytimg.com/vi/%s/oar2.jpg",
	}
)

func youTubeShortThumbnailURLTemplates() []string {
	return slices.Clone(youTubeShortThumbnailURLTemplateList)
}

func fetchYouTubeShortImage(ctx context.Context, url string, options FetchOptions) (*dtypes.ImageData, error) {
	youtubeShortID, err := youtube.ExtractShortID(url)
	if err != nil {
		return nil, err
	}

	var lastErr error

	urls := youTubeShortThumbnailURLTemplates()
	for _, url := range urls {
		imageURL := fmt.Sprintf(url, youtubeShortID)

		imageData, err := fetch.FetchImage(ctx, imageURL, options)
		if imageData != nil {
			return imageData, nil
		}

		lastErr = err
	}

	return nil, lastErr
}

func fetchYouTubeImage(ctx context.Context, mediaURL string, options FetchOptions) (*dtypes.ImageData, error) {
	youtubeID, err := youtube.ExtractID(mediaURL)
	if err != nil {
		return nil, err
	}

	// https://img.youtube.com/vi/<id>/maxresdefault.jpg		- Maximum Resolution / Max Resolution
	// https://img.youtube.com/vi/<id>/hqdefault.jpg			- High Quality
	// https://img.youtube.com/vi/<id>/mqdefault.jpg			- Medium Quality
	// https://img.youtube.com/vi/<id>/default.jpg				- Default Quality / Standard Thumbnail
	// https://img.youtube.com/vi/<id>/sddefault.jpg			- Standard Definition
	imageURL := fmt.Sprintf("https://img.youtube.com/vi/%s/hqdefault.jpg", youtubeID)
	imageData, err := fetch.FetchImage(ctx, imageURL, options)
	if err != nil {
		return nil, err
	}

	return imageData, nil
}

func FetchYouTubeImages(
	ctx context.Context,
	mediaURL string,
	options FetchOptions,
) ([]*dtypes.ImageData, error) {
	if mediaURL == "" {
		return nil, nil
	}

	var images []*dtypes.ImageData

	imageData, _ := fetchYouTubeImage(ctx, mediaURL, options)
	if imageData != nil {
		images = append(images, imageData)
	}

	imageData, _ = fetchYouTubeShortImage(ctx, mediaURL, options)
	if imageData != nil {
		images = append(images, imageData)
	}

	return images, nil
}
