package downloader

import (
	"context"

	dservices "github.com/neosy/elengrab/internal/domain/services"
)

// GetInfo retrieves and parses video formats for the given URL.
func (d *Downloader) GetInfo(ctx context.Context, url string, useCookies bool) (*dservices.DownloaderMediaInfo, error) {
	info, err := d.executor.GetInfo(ctx, url, useCookies)
	if err != nil {
		return nil, err
	}

	return d.mappers.MapMediaInfoToDomain(info), nil
}
