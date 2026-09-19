package channels

import (
	"context"

	"github.com/neosy/elengrab/internal/app/utils/siteimage/internal/fetch"
	"github.com/neosy/elengrab/internal/app/utils/siteimage/internal/types"
	"github.com/neosy/elengrab/internal/app/utils/websource/htmlparser"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
)

func FetchInstagramImages(
	ctx context.Context,
	channelURL string,
	options FetchOptions,
) ([]dtypes.ImageData, error) {
	fetchOptions := types.NewFetchOptions(options)

	body, err := fasthttpx.GetHTML(
		channelURL,
		fasthttpx.ClientOptionWithreadBufferSize(64*1024),
		fasthttpx.ClientOptionWithTimeout(fetchOptions.Timeout),
	)
	if err != nil {
		return nil, errorx.Errorf("failed to get html: %w", err)
	}

	imageURLs, err := htmlparser.ExtractOpenGraphImageURLs(body)
	if err != nil {
		return nil, err
	}

	var loadedImages []dtypes.ImageData

	for _, url := range imageURLs {
		image, err := fetch.FetchImage(ctx, url, fetchOptions)
		if err != nil || image == nil {
			continue
		}

		loadedImages = append(loadedImages, *image)
	}

	return loadedImages, nil
}
