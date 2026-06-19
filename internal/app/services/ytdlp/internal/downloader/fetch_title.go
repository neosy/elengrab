package downloader

import (
	"context"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/helper"
)

func (d *Downloader) FetchTitle(ctx context.Context, url string, options idto.RequestOptions) (string, error) {
	var cookieFileName string
	if options.AllowCookies {
		cookieFileName, _ = helper.CookieFilePathFromURL(url, d.serviceOptions.CookiesDir)
	}
	return d.executor.FetchTitle(ctx, url, idto.WithUseCookies(cookieFileName))
}
