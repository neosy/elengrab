package paths

var (
	PwaManifestPath = pwaFileName("manifest.json").path

	PwaPaths = pwaFileNames{
		"manifest.json",
	}.paths
)

type (
	pwaFileName  string
	pwaFileNames []string
)

func (name pwaFileName) path(dir string) (string, error) {
	return assetPWAPath(string(name), dir)
}

func (names pwaFileNames) paths(dir string) ([]string, error) {
	return assetPWAPaths(names, dir)
}
