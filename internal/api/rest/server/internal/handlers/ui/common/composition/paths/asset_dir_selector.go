package paths

import "github.com/neosy/elengrab/internal/api/rest/server/assets"

type assetDirSelector func(*assets.Assets) string

func assetCssDir(assets *assets.Assets) string {
	return assets.FolderPaths().Css()
}

func assetPwaDir(assets *assets.Assets) string {
	return assets.FolderPaths().Pwa()
}

func assetJsDir(assets *assets.Assets) string {
	return assets.FolderPaths().Js()
}
