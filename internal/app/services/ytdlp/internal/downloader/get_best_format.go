package downloader

import (
	"context"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
)

func (d *Downloader) GetBestFormat(
	ctx context.Context,
	url string,
	format string,
	useCookies bool,
) (*dmedia.MediaInfo, error) {
	info, err := d.executor.GetBestFormat(ctx, url, format, useCookies)
	if err != nil {
		return nil, err
	}

	return d.mappers.MapMediaInfoToDomain(info), nil
}
