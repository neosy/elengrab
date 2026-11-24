package paths

import (
	"context"
	"path/filepath"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
	"github.com/neosy/elengrab/internal/pkg/assetx"
)

var (
	ErrorCssFileName cssFileName = "page-error.css"

	pagesListCssPaths = cssFileNames{
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
		"dialog.css",

		"theme-switcher.css",

		"list-pages.css",

		"menu.css",
		"menu-variants.css",

		"media-result.css",
	}

	indexPageCssPaths = append(
		pagesListCssPaths,

		cssFileNames{
			"index-page.css",
			"upload-media-dialog.css",
			"download-list-dialog.css",
		}...,
	)

	channelPageCssPaths = append(
		pagesListCssPaths,

		cssFileNames{
			"channel-page.css",
		}...,
	)

	userDownloadsPageCssPaths = append(
		pagesListCssPaths,

		cssFileNames{
			"user-downloads-page.css",
			"upload-media-dialog.css",
			"download-list-dialog.css",
		}...,
	)

	adminPageCssPaths = cssFileNames{
		"font-inter.css",
		"base.css",
		"utilities.css",
		"variables.css",

		"components.css",
		"notifications.css",

		"theme-switcher.css",
		"admin-page.css",
	}

	authPageCssPaths = cssFileNames{
		"font-inter.css",
		"base.css",
		"interactions.css",
		"utilities.css",
		"variables.css",

		"theme-switcher.css",
		"auth-page.css",
	}

	pagesDetailCssPaths = cssFileNames{
		"font-inter.css",
		"base.css",
		"interactions.css",
		"utilities.css",
		"variables.css",
		"variables-page.css",

		"components.css",
		"components-header.css",
		"components-footer.css",
		"notifications.css",
		"dialog.css",

		"theme-switcher.css",

		"detail-pages.css",

		"menu.css",
		"menu-variants.css",
	}

	watchPageCssPaths = append(
		pagesDetailCssPaths,

		cssFileNames{
			"watch-page.css",
		}...,
	)

	editMediaPageCssPaths = append(
		pagesDetailCssPaths,

		cssFileNames{
			"edit-media-page.css",
		}...,
	)
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
