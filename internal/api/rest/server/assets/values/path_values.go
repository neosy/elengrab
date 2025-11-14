package avalues

import httppaths "github.com/neosy/elengrab/internal/api/rest/server/paths"

const (
	PathStaticImgKey   = "PathStaticImg"
	PathStaticIconsKey = "PathStaticIcons"
	PathStaticCssKey   = "PathStaticCss"
	PathStaticJsKey    = "PathStaticJs"
)

var PathValues = map[string]string{
	"PathStaticImg":      httppaths.GroupStatic + "/img",
	"PathStaticIcons":    httppaths.GroupStatic + "/img/icons",
	"PathStaticCss":      httppaths.GroupStatic + "/css",
	"PathStaticJs":       httppaths.GroupStatic + "/js",
	"PathDownloaderGrab": httppaths.GroupDownloader + httppaths.PathGrab,
}
