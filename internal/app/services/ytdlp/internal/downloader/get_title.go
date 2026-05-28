package downloader

import (
	"context"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/helper"
)

func (d *Downloader) GetTitle(ctx context.Context, url string, useCookies bool) (string, error) {
	var cookieFileName string
	if useCookies {
		cookieFileName, _ = helper.CookieFilePathFromURL(url, d.serviceOptions.CookiesDir)
	}
	return d.executor.GetTitle(ctx, url, idto.WithUseCookies(cookieFileName))
}
