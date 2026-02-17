package downloader

import (
	"context"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/executor"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
)

func (d *Downloader) GetBestFormat(
	ctx context.Context,
	url string,
	format string,
	useCookies bool,
) (*dmedia.MediaInfo, error) {
	info, err := d.executor.GetBestFormat(
		ctx,
		url,
		format,
		executor.WithUseCookies(useCookies),
	)
	if err != nil {
		return nil, err
	}

	return d.mappers.MapMediaInfoToDomain(info), nil
}
