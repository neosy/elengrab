package avalues

import httppaths "github.com/neosy/elengrab/internal/api/rest/server/paths"

var PathsValues = map[string]string{
	"PathStaticImg":      httppaths.GroupStatic + "/img",
	"PathStaticCss":      httppaths.GroupStatic + "/css",
	"PathStaticJs":       httppaths.GroupStatic + "/js",
	"PathDownloaderGrab": httppaths.GroupDownloader + httppaths.PathGrab,
}
