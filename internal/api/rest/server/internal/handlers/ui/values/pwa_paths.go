package uivalues

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
	path, err := generateHashedFileName(dir, string(name))
	if err != nil {
		return "", err
	}
	return PwaHttpPath(path), nil
}

func (names pwaFileNames) paths(dir string) ([]string, error) {
	paths, err := generateHashedFileNames(dir, names)
	if err != nil {
		return nil, err
	}

	for i, fineName := range paths {
		paths[i] = PwaHttpPath(fineName)
	}

	return paths, nil
}
