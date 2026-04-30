package uivalues

import (
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

const (
	PathItemsHistoryKey = "PathItemsHistory"
)

var (
	PathValues = map[string]any{
		"PathStaticImg":            httppaths.GroupStaticImg,
		"PathStaticIcons":          httppaths.GroupStaticIcon,
		"PathStaticCss":            httppaths.GroupStaticCss,
		"PathStaticJs":             httppaths.GroupStaticJs,
		"PathStaticPwa":            httppaths.GroupStaticPwa,
		"PathStaticThumbnail":      httppaths.GroupStaticThumbnail,
		"PathStaticYoutubeChannel": httppaths.GroupStaticYoutubeChannel,

		"PathAuthRegister": httppaths.GroupAccount + httppaths.PathRegister,
		"PathAuthLogin":    httppaths.GroupAccount + httppaths.PathLogin,

		"PathDownloader":     httppaths.GroupDownloader,
		"PathAccountMenu":    httppaths.GroupDownloader + httppaths.PathAccountMenu,
		"PathRowMenu":        httppaths.GroupDownloader + httppaths.PathFileMenu,
		PathItemsHistoryKey:  httppaths.GroupDownloader + httppaths.PathHistory,
		"PathDownloaderGrab": httppaths.GroupDownloader + httppaths.PathGrab,
		"PathHistorySearch":  httppaths.GroupDownloader + httppaths.PathSearch,
	}
)

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

func ThumbnailHttpPath(id string) string {
	return httppaths.GroupStaticThumbnail + "/" + id
}

func YoutubeChannelHttpPath(channelId string) string {
	return httppaths.GroupStaticYoutubeChannel + "/" + channelId
}
