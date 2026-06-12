package paths

import (
	"path/filepath"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
)

var (
	assetCssPaths = newAssetPaths(CssPath)

	assetPWAPath  = newAssetPath(PwaPath)
	assetPWAPaths = newAssetPaths(PwaPath)

	assetJsNameKeyPath = newAssetNameKeyPaths(JsPath)
)

type (
	assetPath         func(fileName, dir string) (string, error)
	assetPaths        func(fileNames []string, dir string) ([]string, error)
	assetNameKeyPaths func(fileNames []string, dir string) (map[string]string, error)
)

func newAssetPath(httpPath func(fileName string) string) assetPath {
	return func(fileName, dir string) (string, error) {
		filePath := filepath.Join(dir, fileName)
		file, err := assets.ReadAssetFileWithHash(filePath)
		if err != nil {
			return "", err
		}
		return httpPath(file.FileNameWithHash()), nil
	}
}

func newAssetPaths(httpPath func(fileName string) string) assetPaths {
	return func(fileNames []string, dir string) ([]string, error) {
		paths, err := fileNamesWithHash(dir, fileNames)
		if err != nil {
			return nil, err
		}

		for i, fineName := range paths {
			paths[i] = httpPath(fineName)
		}

		return paths, nil
	}
}

func newAssetNameKeyPaths(httpPath func(fileName string) string) assetNameKeyPaths {
	return func(fileNames []string, dir string) (map[string]string, error) {
		paths, err := fileNamesWithHash(dir, fileNames)
		if err != nil {
			return nil, err
		}

		var jsMap = make(jsImportMap, len(paths))
		for i, newName := range paths {
			key := JsPath(fileNames[i])
			jsMap[key] = httpPath(newName)
		}

		return jsMap, nil
	}
}
