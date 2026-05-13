package uivalues

import (
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

const (
	PathDownloaderHistoryKey = "PathDownloaderHistory"
)

type (
	basePaths struct {
		StaticImg            string
		StaticIcons          string
		StaticCss            string
		StaticJs             string
		StaticPwa            string
		StaticThumbnail      string
		StaticYoutubeChannel string

		AuthRegister string
		AuthLogin    string

		Downloader        string
		AccountMenu       string
		RowMenu           string
		DownloaderHistory string
		DownloaderGrab    string
		DownloaderSearch  string
		DownloaderEvents  string
	}
)

var (
	basePathsDefault = basePaths{
		StaticImg:            httppaths.GroupStaticImg,
		StaticIcons:          httppaths.GroupStaticIcon,
		StaticCss:            httppaths.GroupStaticCss,
		StaticJs:             httppaths.GroupStaticJs,
		StaticPwa:            httppaths.GroupStaticPwa,
		StaticThumbnail:      httppaths.GroupStaticThumbnail,
		StaticYoutubeChannel: httppaths.GroupStaticYoutubeChannel,

		AuthRegister: httppaths.GroupAccount + httppaths.PathRegister,
		AuthLogin:    httppaths.GroupAccount + httppaths.PathLogin,

		Downloader:        httppaths.GroupDownloader,
		AccountMenu:       httppaths.GroupDownloader + httppaths.PathAccountMenu,
		RowMenu:           httppaths.GroupDownloader + httppaths.PathMediaItemMenu,
		DownloaderHistory: httppaths.GroupDownloader + httppaths.PathHistory,
		DownloaderGrab:    httppaths.GroupDownloader + httppaths.PathGrab,
		DownloaderSearch:  httppaths.GroupDownloader + httppaths.PathSearch,
		DownloaderEvents:  httppaths.GroupDownloader + httppaths.PathEvents,
	}
)

func NewBasePaths() basePaths {
	return basePathsDefault
}

func ImageHttpPath(fileName string) string {
	return httppaths.GroupStaticImg + "/" + fileName
}

func IconHttpPath(fileName string) string {
	return httppaths.GroupStaticIcon + "/" + fileName
}

func CssHttpPath(fileName string) string {
	return httppaths.GroupStaticCss + "/" + fileName
}

func JsHttpPath(fileName string) string {
	return httppaths.GroupStaticJs + "/" + fileName
}

func PwaHttpPath(fileName string) string {
	return httppaths.GroupStaticPwa + "/" + fileName
}

func ThumbnailHttpPath(id string) string {
	return httppaths.GroupStaticThumbnail + "/" + id
}

func YoutubeChannelHttpPath(channelId string) string {
	return httppaths.GroupStaticYoutubeChannel + "/" + channelId
}
