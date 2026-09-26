package paths

import (
	"context"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
)

type AssetPaths struct {
	loaders assetPathLoaders

	IndexPageCssPaths   func() ([]string, error)
	ChannelPageCssPaths func() ([]string, error)

	AuthPageCssPaths      func() ([]string, error)
	AdminPageCssPaths     func() ([]string, error)
	WatchPageCssPaths     func() ([]string, error)
	EditMediaPageCssPaths func() ([]string, error)

	IndexPageJsPaths   func(legacy bool) ([]JsScript, error)
	ChannelPageJsPaths func(legacy bool) ([]JsScript, error)

	AuthPageJsPaths      func(legacy bool) ([]JsScript, error)
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

		IndexPageCssPaths:   indexPageCssPaths.newLoader(assets, loaders.cssPaths),
		ChannelPageCssPaths: channelPageCssPaths.newLoader(assets, loaders.cssPaths),

		AuthPageCssPaths:      authPageCssPaths.newLoader(assets, loaders.cssPaths),
		AdminPageCssPaths:     adminPageCssPaths.newLoader(assets, loaders.cssPaths),
		WatchPageCssPaths:     watchPageCssPaths.newLoader(assets, loaders.cssPaths),
		EditMediaPageCssPaths: editMediaPageCssPaths.newLoader(assets, loaders.cssPaths),

		IndexPageJsPaths:   indexPageJsPaths.newLoader(ctx, assets),
		ChannelPageJsPaths: channelPageJsPaths.newLoader(ctx, assets),

		AuthPageJsPaths:      authPageJsPaths.newLoader(ctx, assets),
		AdminPageJsPaths:     adminPageJsPaths.newLoader(ctx, assets),
		WatchPageJsPaths:     watchPageJsPaths.newLoader(ctx, assets),
		EditMediaPageJsPaths: editMediaPageJsPaths.newLoader(ctx, assets),

		PwaManifestPath: pwaManifestPath.newLoader(assets, loaders.pwaPath),
		PwaPaths:        pwaPaths.newLoader(assets, loaders.pwaPaths),
	}
}
