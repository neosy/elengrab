package helper

import (
	"context"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	"github.com/neosy/elengrab/internal/pkg/imgx"
)

// FetchImage downloads an image from the given URL and returns it.
func FetchImage(ctx context.Context, imgURL string, options idto.RequestOptions) (*dtypes.ImageData, error) {
	raw, format, err := httpx.FetchImage(
		ctx,
		imgURL,
		httpx.MethodGetOptions{Limit: options.LimitBytes},
		httpx.ClientOptionWithTimeout(options.Timeout),
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
