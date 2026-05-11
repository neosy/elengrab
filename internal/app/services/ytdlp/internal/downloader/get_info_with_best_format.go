package downloader

import (
	"context"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/executor"
	dservices "github.com/neosy/elengrab/internal/domain/services"
)

func (d *Downloader) GetInfoWithBestFormat(
	ctx context.Context,
	url string,
	format string,
	useCookies bool,
) (*dservices.DownloaderMediaInfo, error) {
	info, err := d.executor.GetInfoWithBestFormat(
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
