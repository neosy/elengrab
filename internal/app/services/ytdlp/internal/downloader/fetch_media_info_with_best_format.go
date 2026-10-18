package downloader

import (
	"context"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/helper"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

func (d *Downloader) FetchInfoWithBestFormat(
	ctx context.Context,
	mediaURL string,
	format string,
	options idto.RequestOptions,
) (*pservices.DownloaderMeta, error) {
	var cookieFilePath string
	if options.AllowCookies && d.serviceOptions.AllowCookies {
		cookieFilePath, _ = helper.CookieFilePathFromURL(mediaURL, d.serviceOptions.CookiesDir)
	}

	info, err := d.executor.FetchInfoWithBestFormat(
		ctx,
		mediaURL,
		format,
		idto.WithUseCookies(cookieFilePath),
	)
	if err != nil {
		return nil, err
	}

	return d.mappers.MapExtractInfoToDownloaderMeta(info), nil
}
