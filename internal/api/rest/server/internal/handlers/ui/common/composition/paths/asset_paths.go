package paths

import (
	"context"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
)

type AssetPaths struct {
	loaders assetPathLoaders

	AuthPageCssPaths      func() ([]string, error)
	IndexPageCssPaths     func() ([]string, error)
	AdminPageCssPaths     func() ([]string, error)
	WatchPageCssPaths     func() ([]string, error)
	EditMediaPageCssPaths func() ([]string, error)

	AuthPageJsPaths      func(legacy bool) ([]JsScript, error)
	IndexPageJsPaths     func(legacy bool) ([]JsScript, error)
	AdminPageJsPaths     func(legacy bool) ([]JsScript, error)
	WatchPageJsPaths     func(legacy bool) ([]JsScript, error)
	EditMediaPageJsPaths func(legacy bool) ([]JsScript, error)

	PwaManifestPath func() (string, error)
	PwaPaths        func() ([]string, error)
}

func NewAssetPaths(assets *assets.Assets) AssetPaths {
	ctx := context.Background()
	loaders := newAssetPathLoaders(ctx)

	return AssetPaths{
		loaders: loaders,

		AuthPageCssPaths:      authPageCssPaths.newLoader(assets, loaders.cssPaths),
		IndexPageCssPaths:     indexPageCssPaths.newLoader(assets, loaders.cssPaths),
		AdminPageCssPaths:     adminPageCssPaths.newLoader(assets, loaders.cssPaths),
		WatchPageCssPaths:     watchPageCssPaths.newLoader(assets, loaders.cssPaths),
		EditMediaPageCssPaths: editMediaPageCssPaths.newLoader(assets, loaders.cssPaths),

		AuthPageJsPaths:      authPageJsPaths.newLoader(ctx, assets),
		IndexPageJsPaths:     indexPageJsPaths.newLoader(ctx, assets),
		AdminPageJsPaths:     adminPageJsPaths.newLoader(ctx, assets),
		WatchPageJsPaths:     watchPageJsPaths.newLoader(ctx, assets),
		EditMediaPageJsPaths: editMediaPageJsPaths.newLoader(ctx, assets),

		PwaManifestPath: pwaManifestPath.newLoader(assets, loaders.pwaPath),
		PwaPaths:        pwaPaths.newLoader(assets, loaders.pwaPaths),
	}
}
