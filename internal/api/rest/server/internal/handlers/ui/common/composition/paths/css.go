package paths

import (
	"context"
	"path/filepath"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
	"github.com/neosy/elengrab/internal/pkg/assetx"
)

var (
	ErrorCssFileName cssFileName = "page-error.css"

	indexPageCssPaths = cssFileNames{
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
		"video-preview.css",
		"notifications.css",
	}

	adminPageCssPaths = cssFileNames{
		"font-inter.css",
		"base.css",
		"utilities.css",
		"variables.css",
		"notifications.css",

		"page-admin.css",
		"theme-switcher.css",
	}

	authPageCssPaths = cssFileNames{
		"font-inter.css",
		"base.css",
		"interactions.css",
		"utilities.css",
		"variables.css",

		"page-auth.css",
		"theme-switcher.css",
	}

	watchPageCssPaths = cssFileNames{
		"font-inter.css",
		"base.css",
		"interactions.css",
		"utilities.css",
		"variables.css",

		"page-watch.css",
		"theme-switcher.css",

		"notifications.css",
	}

	editMediaPageCssPaths = cssFileNames{
		"font-inter.css",
		"base.css",
		"interactions.css",
		"utilities.css",
		"variables.css",
		"variables-page.css",

		"page-edit-media.css",
		"theme-switcher.css",

		"notifications.css",
	}
)

type (
	cssFileName  string
	cssFileNames []string
)

func (names cssFileNames) newLoader(assets *assets.Assets, loader loaderAssetPaths) func() ([]string, error) {
	return func() ([]string, error) {
		return loader(names, assets)
	}
}

func (name cssFileName) Raw(ctx context.Context, assets *assets.Assets) ([]byte, error) {
	filePath := filepath.Join(assets.FolderPaths().Css(), string(name))

	file, err := assets.ReadAssetFile(ctx, filePath)
	if err != nil {
		return nil, err
	}

	return assetx.MinifyCSS(file.Raw), nil
}
