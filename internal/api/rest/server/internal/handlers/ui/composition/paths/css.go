package paths

import (
	"os"
	"path/filepath"

	"github.com/neosy/elengrab/internal/pkg/assetx"
)

var (
	ErrorCssFileName cssFileName = "page-error.css"

	IndexCssPaths = cssFileNames{
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

	AdminCssPaths = cssFileNames{
		"font-inter.css",
		"base.css",
		"utilities.css",
		"variables.css",
		"notifications.css",

		"page-admin.css",
		"theme-switcher.css",
	}.paths

	AuthCssPaths = cssFileNames{
		"font-inter.css",
		"base.css",
		"interactions.css",
		"utilities.css",
		"variables.css",

		"page-auth.css",
		"theme-switcher.css",
	}.paths

	WatchCssPaths = cssFileNames{
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
	filePath := filepath.Join(dir, string(name))

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return assetx.MinifyCSS(data), nil
}
