package paths

import (
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

const (
	PathDownloaderHistoryKey = "PathDownloaderHistory"
)

type (
	HttpPaths struct {
		StaticImages          string
		StaticIcons           string
		StaticCss             string
		StaticJs              string
		StaticPwa             string
		StaticThumbnails      string
		StaticYoutubeChannels string

		AuthRegister string
		AuthLogin    string

		Downloader        string
		AccountMenu       string
		RowMenu           string
		SettingsMenu      string
		DownloaderHistory string
		DownloaderGrab    string
		DownloaderSearch  string
		DownloaderEvents  string
	}
)

var (
	httpPathsDefault = HttpPaths{
		StaticImages:          httppaths.GroupStaticImages,
		StaticIcons:           httppaths.GroupStaticIcons,
		StaticCss:             httppaths.GroupStaticCss,
		StaticJs:              httppaths.GroupStaticJs,
		StaticPwa:             httppaths.GroupStaticPwa,
		StaticThumbnails:      httppaths.GroupStaticThumbnails,
		StaticYoutubeChannels: httppaths.GroupStaticYoutubeChannels,

		AuthRegister: httppaths.GroupAccount + httppaths.PathRegister,
		AuthLogin:    httppaths.GroupAccount + httppaths.PathLogin,

		Downloader:        httppaths.GroupDownloader,
		AccountMenu:       httppaths.GroupDownloader + httppaths.PathAccountMenu,
		SettingsMenu:      httppaths.GroupDownloader + httppaths.PathSettingsMenu,
		RowMenu:           httppaths.GroupDownloader + httppaths.PathMediaItemMenu,
		DownloaderHistory: httppaths.GroupDownloader + httppaths.PathHistory,
		DownloaderGrab:    httppaths.GroupDownloader + httppaths.PathGrab,
		DownloaderSearch:  httppaths.GroupDownloader + httppaths.PathSearch,
		DownloaderEvents:  httppaths.GroupDownloader + httppaths.PathEvents,
	}
)

func NewHttpPaths() HttpPaths {
	return httpPathsDefault
}

func ImagePath(fileName string) string {
	return httppaths.GroupStaticImages + "/" + fileName
}

func IconPath(fileName string) string {
	return httppaths.GroupStaticIcons + "/" + fileName
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
	return httppaths.GroupStaticThumbnails + "/" + id
}

func YoutubeChannelPath(channelId string) string {
	return httppaths.GroupStaticYoutubeChannels + "/" + channelId
}
