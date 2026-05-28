package downloader

import (
	"context"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/helper"
	dservices "github.com/neosy/elengrab/internal/domain/services"
)

// GetInfo retrieves and parses video formats for the given URL.
func (d *Downloader) GetInfo(
	ctx context.Context,
	url string,
	useCookies bool,
) (*dservices.DownloaderMediaInfo, error) {
	var cookieFileName string
	if useCookies {
		cookieFileName, _ = helper.CookieFilePathFromURL(url, d.serviceOptions.CookiesDir)
	}

	info, err := d.executor.GetInfo(ctx, url, idto.WithUseCookies(cookieFileName))
	if err != nil {
		return nil, err
	}

	return d.mappers.MapMediaInfoToDomain(info), nil
}
