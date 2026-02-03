package logofetcher

import (
	"context"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
	"github.com/neosy/elengrab/pkg/httpx"
)

func (lf *SiteLogoFetcher) downloadImage(ctx context.Context, imgURL string) (*dmedia.ImageData, error) {
	data, format, err := httpx.GetImage(
		ctx,
		imgURL,
		httpx.GetOptions{Limit: limitImage},
		httpx.ClientOptionWithTimeout(downloadImageTimeout),
	)
	if err != nil {
		return nil, err
	}
	return &dmedia.ImageData{
		URL:    imgURL,
		Raw:    data,
		Format: format,
	}, nil
}
