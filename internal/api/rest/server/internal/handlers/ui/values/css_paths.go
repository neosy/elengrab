package uivalues

import (
	"github.com/neosy/elengrab/internal/pkg/httpx"
)

var (
	CssErrorFileName cssFileName = "error.css"
)

var (
	CssIndexPaths = cssFileNames{
		"font-inter.css",
		"base.css",
		"interactions.css",
		"utilities.css",
		"variables.css",

		"page-index.css",
		"theme-switcher.css",

		"menu.css",
		"menu-variants.css",

		"grab-form.css",
		"result-rows.css",

		"player.css",
		"notifications.css",
	}.paths

	CssAuthPaths = cssFileNames{
		"font-inter.css",
		"base.css",
		"interactions.css",
		"utilities.css",
		"variables.css",

		"page-auth.css",
		"theme-switcher.css",
	}.paths

	CssWatchPaths = cssFileNames{
		"font-inter.css",
		"base.css",
		"interactions.css",
		"utilities.css",
		"variables.css",

		"page-watch.css",
		"theme-switcher.css",
	}.paths
)

type (
	cssFileName  string
	cssFileNames []string
)

func (names cssFileNames) paths(dir string) ([]string, error) {
	return assetCssPaths(names, dir)
}

func (name cssFileName) Raw(dir string) ([]byte, error) {
	data, err := fileRaw(string(name), dir)
	if err != nil {
		return nil, err
	}
	return httpx.MinifyCSS(data), nil
}
