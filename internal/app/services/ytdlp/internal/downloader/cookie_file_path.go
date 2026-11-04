package downloader

import "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/helper"

func (d *Downloader) cookieFilePath(url string) (string, error) {
	if !d.serviceOptions.AllowCookies {
		return "", nil
	}

	cookieFilePath, err := helper.CookieFilePathFromURL(url, d.serviceOptions.CookiesDir)
	if err != nil {
		return "", err
	}

	return cookieFilePath, nil
}
