package paths

import "github.com/neosy/elengrab/internal/api/rest/server/assets"

var (
	pwaManifestPath = pwaFileName("manifest.json")

	pwaPaths = pwaFileNames{
		"manifest.json",
	}
)

type (
	pwaFileName  string
	pwaFileNames []string
)

func (name pwaFileName) newLoader(assets *assets.Assets, loader loaderAssetPath) func() (string, error) {
	return func() (string, error) {
		return loader(string(name), assets)
	}
}

func (names pwaFileNames) newLoader(assets *assets.Assets, loader loaderAssetPaths) func() ([]string, error) {
	return func() ([]string, error) {
		return loader(names, assets)
	}
}
