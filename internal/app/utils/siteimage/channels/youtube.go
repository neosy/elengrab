package channels

import (
	"context"
	"encoding/json"

	"github.com/neosy/elengrab/internal/app/utils/siteimage/internal/fetch"
	"github.com/neosy/elengrab/internal/app/utils/siteimage/internal/types"
	"github.com/neosy/elengrab/internal/app/utils/websource/htmlparser"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
)

// FetchYouTubeImages fetches the HTML of a YouTube channel page,
// extracts the avatar JSON block, and returns all image sources.
// url: full URL of the YouTube channel
func FetchYouTubeImages(ctx context.Context, channelURL string, options FetchOptions) ([]dtypes.ImageData, error) {
	fetchOptions := types.NewFetchOptions(options)

	body, err := fasthttpx.GetHTML(
		channelURL,
		fasthttpx.ClientOptionWithreadBufferSize(64*1024),
		fasthttpx.ClientOptionWithTimeout(fetchOptions.Timeout),
	)
	if err != nil {
		return nil, errorx.Errorf("failed to get html: %w", err)
	}

	// Full key path to the avatar sources inside decoratedAvatarViewModel
	fullKey := `"decoratedAvatarViewModel":{"avatar":{"avatarViewModel":{"image":{"sources":[`
	jsonArray, err := htmlparser.ExtractJSONArray(body, fullKey)
	if err != nil {
		return nil, err
	}

	var sources []ImageSource
	var loadedImages []dtypes.ImageData
	if err := json.Unmarshal([]byte(jsonArray), &sources); err != nil {
		return nil, err
	}

	for _, src := range sources {
		image, err := fetch.FetchImage(ctx, src.URL, fetchOptions)
		if err != nil || image == nil {
			continue
		}

		loadedImages = append(loadedImages, *image)
	}

	return loadedImages, nil
}
