package ucdownloader

import (
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *YouTubeDownloader) Download(url string, options *ddownload.DownloadOptions) (*ddownload.DownloadResponse, error) {
	resp, err := uc.downloaderSrv.Download(
		url,
		options,
	)
	return resp, err
}
