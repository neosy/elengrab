package downloader

import (
	"context"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/helper"
	dservices "github.com/neosy/elengrab/internal/domain/services"
)

func (d *Downloader) FetchInfoWithBestFormat(
	ctx context.Context,
	url string,
	format string,
	options idto.RequestOptions,
) (*dservices.DownloaderMediaInfo, error) {
	var cookieFileName string
	if options.AllowCookies {
		cookieFileName, _ = helper.CookieFilePathFromURL(url, d.serviceOptions.CookiesDir)
	}
	info, err := d.executor.FetchInfoWithBestFormat(
		ctx,
		url,
		format,
		idto.WithUseCookies(cookieFileName),
	)
	if err != nil {
		return nil, err
	}

	return d.mappers.MapMediaInfoToDomain(info), nil
}
