package downloader

import (
	"context"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
)

// GetFormats retrieves and parses video formats for the given URL.
func (d *Downloader) GetFormats(ctx context.Context, url string, useCookies bool) (*dmedia.MediaInfo, error) {
	info, err := d.executor.GetFormats(ctx, url, useCookies)
	if err != nil {
		return nil, err
	}

	return d.mappers.MapMediaInfoToDomain(info), nil
}
