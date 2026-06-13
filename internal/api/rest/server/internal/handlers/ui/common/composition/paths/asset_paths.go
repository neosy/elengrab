package paths

import "github.com/neosy/elengrab/internal/api/rest/server/assets"

type AssetPaths struct {
	loaders assetPathLoaders

	AuthPageCssPaths  func() ([]string, error)
	IndexPageCssPaths func() ([]string, error)
	AdminPageCssPaths func() ([]string, error)
	WatchPageCssPaths func() ([]string, error)

	AuthPageJsPaths  func() ([]JsScript, error)
	IndexPageJsPaths func() ([]JsScript, error)
	AdminPageJsPaths func() ([]JsScript, error)
	WatchPageJsPaths func() ([]JsScript, error)

	PwaManifestPath func() (string, error)
	PwaPaths        func() ([]string, error)
}

func NewAssetPaths(assets *assets.Assets) AssetPaths {
	loaders := newAssetPathLoaders()
	return AssetPaths{
		loaders: loaders,

		AuthPageCssPaths:  authPageCssPaths.newLoader(assets, loaders.cssPaths),
		IndexPageCssPaths: indexPageCssPaths.newLoader(assets, loaders.cssPaths),
		AdminPageCssPaths: adminPageCssPaths.newLoader(assets, loaders.cssPaths),
		WatchPageCssPaths: watchPageCssPaths.newLoader(assets, loaders.cssPaths),

		AuthPageJsPaths:  authPageJsPaths.newLoader(assets),
		IndexPageJsPaths: indexPageJsPaths.newLoader(assets),
		AdminPageJsPaths: adminPageJsPaths.newLoader(assets),
		WatchPageJsPaths: watchPageJsPaths.newLoader(assets),

		PwaManifestPath: pwaManifestPath.newLoader(assets, loaders.pwaPath),
		PwaPaths:        pwaPaths.newLoader(assets, loaders.pwaPaths),
	}
}
