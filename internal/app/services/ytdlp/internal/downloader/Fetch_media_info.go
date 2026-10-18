package downloader

import (
	"context"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/helper"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

// GetInfo retrieves and parses video formats for the given URL.
func (d *Downloader) FetchInfo(
	ctx context.Context,
	url string,
	options idto.RequestOptions,
) (*pservices.DownloaderMeta, error) {
	var cookieFilePath string
	if options.AllowCookies {
		cookieFilePath, _ = helper.CookieFilePathFromURL(url, d.serviceOptions.CookiesDir)
	}

	info, err := d.executor.FetchInfo(ctx, url, idto.WithUseCookies(cookieFilePath))
	if err != nil {
		return nil, err
	}

	return d.mappers.MapExtractInfoToDownloaderMeta(info), nil
}
