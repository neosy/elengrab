package fetch

import (
	"context"

	"github.com/neosy/elengrab/internal/app/utils/siteimage/internal/types"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	"github.com/neosy/elengrab/internal/pkg/imgx"
)

// FetchImage downloads an image from the given URL and returns it.
func FetchImage(ctx context.Context, imgURL string, options Options) (*dtypes.ImageData, error) {
	imageOptions := types.NewFetchOptions(options)

	raw, format, err := httpx.FetchImage(
		ctx,
		imgURL,
		httpx.MethodGetOptions{Limit: imageOptions.LimitBytes},
		httpx.ClientOptionWithTimeout(imageOptions.Timeout),
	)
	if err != nil {
		return nil, err
	}

	imageFormat, err := dtypes.ParseImageFormat(format)
	if err != nil {
		return nil, err
	}

	size, _ := imgx.ImageSize(raw)

	return &dtypes.ImageData{
		URL:    imgURL,
		Format: imageFormat,

		Width:  size.Width,
		Height: size.Height,

		Raw: raw,
	}, nil
}
