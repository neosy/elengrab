package downloader

import (
	"context"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/helper"
	dservices "github.com/neosy/elengrab/internal/domain/services"
)

// GetInfo retrieves and parses video formats for the given URL.
func (d *Downloader) FetchInfo(
	ctx context.Context,
	url string,
	options idto.RequestOptions,
) (*dservices.DownloaderMediaInfo, error) {
	var cookieFilePath string
	if options.AllowCookies {
		cookieFilePath, _ = helper.CookieFilePathFromURL(url, d.serviceOptions.CookiesDir)
	}

	info, err := d.executor.FetchInfo(ctx, url, idto.WithUseCookies(cookieFilePath))
	if err != nil {
		return nil, err
	}

	return d.mappers.MapMediaInfoToDomain(info), nil
}
