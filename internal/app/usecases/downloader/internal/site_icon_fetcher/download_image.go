package iconfetcher

import (
	"context"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/httpx"
)

// fetchImage downloads an image from the given URL and returns it.
func (lf *SiteIconFetcher) fetchImage(ctx context.Context, imgURL string) (*dtypes.ImageData, error) {
	data, format, err := httpx.FetchImage(
		ctx,
		imgURL,
		httpx.MethodGetOptions{Limit: limitImage},
		httpx.ClientOptionWithTimeout(downloadImageTimeout),
	)
	if err != nil {
		return nil, err
	}

	imageFormat, err := dtypes.ParseImageFormat(format)
	if err != nil {
		lf.logger.Warn(
			"Failed to parse image format",
			"format", format,
			"error", err,
		)
		return nil, err
	}

	return &dtypes.ImageData{
		URL:    imgURL,
		Raw:    data,
		Format: imageFormat,
	}, nil
}
