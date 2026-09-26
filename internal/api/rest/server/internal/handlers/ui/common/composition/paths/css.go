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

		"components.css",
		"components-header.css",
		"components-footer.css",
		"player.css",
		"video-preview.css",
		"notifications.css",

		"theme-switcher.css",

		"pages-list.css",
		"page-index.css",

		"menu.css",
		"menu-variants.css",

		"grab-form.css",
		"media-result.css",
	}

	adminPageCssPaths = cssFileNames{
		"font-inter.css",
		"base.css",
		"utilities.css",
		"variables.css",

		"components.css",
		"notifications.css",

		"theme-switcher.css",
		"page-admin.css",
	}

	authPageCssPaths = cssFileNames{
		"font-inter.css",
		"base.css",
		"interactions.css",
		"utilities.css",
		"variables.css",

		"theme-switcher.css",
		"page-auth.css",
	}

	watchPageCssPaths = cssFileNames{
		"font-inter.css",
		"base.css",
		"interactions.css",
		"utilities.css",
		"variables.css",

		"components.css",
		"notifications.css",

		"theme-switcher.css",
		"page-watch.css",
	}

	editMediaPageCssPaths = cssFileNames{
		"font-inter.css",
		"base.css",
		"interactions.css",
		"utilities.css",
		"variables.css",
		"variables-page.css",

		"components.css",
		"notifications.css",

		"theme-switcher.css",
		"page-edit-media.css",
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
