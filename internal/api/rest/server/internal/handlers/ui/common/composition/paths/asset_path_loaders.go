package paths

import (
	"path/filepath"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
)

type (
	loaderAssetPath         func(fileName string, assets *assets.Assets) (string, error)
	loaderAssetPaths        func(fileNames []string, assets *assets.Assets) ([]string, error)
	loaderAssetNameKeyPaths func(fileNames []string, assets *assets.Assets) (map[string]string, error)

	assetPathLoaders struct {
		cssPaths loaderAssetPaths

		pwaPath  loaderAssetPath
		pwaPaths loaderAssetPaths

		jsNameKeyPath loaderAssetNameKeyPaths
	}
)

func newAssetPathLoaders() assetPathLoaders {
	return assetPathLoaders{
		cssPaths: newLoaderAssetPaths(assetCssDir, CssPath),

		pwaPath:  newLoaderAssetPath(assetPwaDir, PwaPath),
		pwaPaths: newLoaderAssetPaths(assetPwaDir, PwaPath),

		jsNameKeyPath: newLoaderAssetNameKeyPaths(assetJsDir, JsPath),
	}
}

func newLoaderAssetPath(dirSelector assetDirSelector, httpPath func(fileName string) string) loaderAssetPath {
	return func(fileName string, assets *assets.Assets) (string, error) {
		filePath := filepath.Join(dirSelector(assets), fileName)
		file, err := assets.ReadAssetFile(filePath)
		if err != nil {
			return "", err
		}
		return httpPath(file.FileNameWithHash()), nil
	}
}

func newLoaderAssetPaths(dirSelector assetDirSelector, httpPath func(fileName string) string) loaderAssetPaths {
	return func(fileNames []string, assets *assets.Assets) ([]string, error) {
		paths, err := loadFileNamesWithHash(fileNames, dirSelector, assets)
		if err != nil {
			return nil, err
		}

		for i, fineName := range paths {
			paths[i] = httpPath(fineName)
		}

		return paths, nil
	}
}

func newLoaderAssetNameKeyPaths(dirSelector assetDirSelector, httpPath func(fileName string) string) loaderAssetNameKeyPaths {
	return func(fileNames []string, assets *assets.Assets) (map[string]string, error) {
		paths, err := loadFileNamesWithHash(fileNames, dirSelector, assets)
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
