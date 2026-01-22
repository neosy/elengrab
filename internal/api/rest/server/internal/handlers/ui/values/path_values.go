package uivalues

import httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"

const (
	PathFileRowKey        = "PathFileRow"
	PathFilesHistoryKey   = "PathFilesHistory"
	PathStaticImgKey      = "PathStaticImg"
	PathStaticIconsKey    = "PathStaticIcons"
	PathStaticCssKey      = "PathStaticCss"
	PathStaticJsKey       = "PathStaticJs"
	PathStaticPwaKey      = "PathStaticPwa"
	PathItemsHistoryKey   = "PathItemsHistory"
	PathDownloaderGrabKey = "PathDownloaderGrab"
)

var PathValues = map[string]any{
	PathStaticImgKey:      httppaths.GroupStatic + "/img",
	PathStaticIconsKey:    httppaths.GroupStatic + "/img/icons",
	PathStaticCssKey:      httppaths.GroupStatic + "/css",
	PathStaticJsKey:       httppaths.GroupStatic + "/js",
	PathStaticPwaKey:      httppaths.GroupStatic + "/pwa",
	PathItemsHistoryKey:   httppaths.GroupDownloader + httppaths.PathHistory,
	PathDownloaderGrabKey: httppaths.GroupDownloader + httppaths.PathGrab,
}
