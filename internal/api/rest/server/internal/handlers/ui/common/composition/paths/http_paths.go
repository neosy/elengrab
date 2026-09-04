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

		Downloader       string
		AccountMenu      string
		RowMenu          string
		SettingsMenu     string
		DownloaderItems  string
		DownloaderGrab   string
		DownloaderSearch string
		DownloaderEvents string
	}
)

var (
	httpPathsDefault = HttpPaths{
		StaticImages:          httppaths.StaticImagesGroup,
		StaticIcons:           httppaths.StaticIconsGroup,
		StaticCss:             httppaths.StaticCssGroup,
		StaticJs:              httppaths.StaticJsGroup,
		StaticPwa:             httppaths.StaticPwaGroup,
		StaticThumbnails:      httppaths.StaticThumbnailsGroup,
		StaticYoutubeChannels: httppaths.StaticYoutubeChannelsGroup,

		AuthRegister: httppaths.AuthRegisterPath,
		AuthLogin:    httppaths.AuthLoginPath,

		Downloader:       httppaths.DownloaderGroup,
		AccountMenu:      httppaths.DownloaderAccountMenuPath,
		SettingsMenu:     httppaths.DownloaderSettingsMenuPath,
		RowMenu:          httppaths.DownloaderGroup + httppaths.MediaItemMenuPath,
		DownloaderItems:  httppaths.DownloaderItemsPath,
		DownloaderGrab:   httppaths.DownloaderGrabPath,
		DownloaderSearch: httppaths.DownloaderSearchPath,
		DownloaderEvents: httppaths.DownloaderEventsPath,
	}
)

func NewHttpPaths() HttpPaths {
	return httpPathsDefault
}

func ImagePath(fileName string) string {
	return httppaths.StaticImagesGroup + "/" + fileName
}

func IconPath(fileName string) string {
	return httppaths.StaticIconsGroup + "/" + fileName
}

func CssPath(fileName string) string {
	return httppaths.StaticCssGroup + "/" + fileName
}

func JsPath(fileName string) string {
	return httppaths.StaticJsGroup + "/" + fileName
}

func PwaPath(fileName string) string {
	return httppaths.StaticPwaGroup + "/" + fileName
}

func ThumbnailPath(id string) string {
	return httppaths.StaticThumbnailsGroup + "/" + id
}

func YoutubeChannelPath(channelId string) string {
	return httppaths.StaticYoutubeChannelsGroup + "/" + channelId
}
