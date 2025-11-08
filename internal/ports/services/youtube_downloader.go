package pservices

import (
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dyoutubeinfo "github.com/neosy/elengrab/internal/domain/youtube_info"
)

type YouTubeDownloader interface {
	GetTitle(url string) (string, error)
	GetFormats(url string) (*dyoutubeinfo.YouTubeInfo, error)
	GetBestFormat(url string) (*dyoutubeinfo.YouTubeInfo, error)
	Download(url string, options *ddownload.DownloadOptions) (*ddownload.DownloadResponse, error)
}
