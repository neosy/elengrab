package paths

import (
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

const (
	PathDownloaderHistoryKey = "PathDownloaderHistory"
)

type (
	Paths struct {
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
	basePathsDefault = Paths{
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

func NewPaths() Paths {
	return basePathsDefault
}

func ImagePath(fileName string) string {
	return httppaths.GroupStaticImg + "/" + fileName
}

func IconPath(fileName string) string {
	return httppaths.GroupStaticIcon + "/" + fileName
}

func CssPath(fileName string) string {
	return httppaths.GroupStaticCss + "/" + fileName
}

func JsPath(fileName string) string {
	return httppaths.GroupStaticJs + "/" + fileName
}

func PwaPath(fileName string) string {
	return httppaths.GroupStaticPwa + "/" + fileName
}

func ThumbnailPath(id string) string {
	return httppaths.GroupStaticThumbnail + "/" + id
}

func YoutubeChannelPath(channelId string) string {
	return httppaths.GroupStaticYoutubeChannel + "/" + channelId
}
