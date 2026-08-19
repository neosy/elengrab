package paths

import (
	"context"
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

func newAssetPathLoaders(ctx context.Context) assetPathLoaders {
	return assetPathLoaders{
		cssPaths: newLoaderAssetPaths(ctx, assetCssDir, CssPath),

		pwaPath:  newLoaderAssetPath(ctx, assetPwaDir, PwaPath),
		pwaPaths: newLoaderAssetPaths(ctx, assetPwaDir, PwaPath),

		jsNameKeyPath: newLoaderAssetNameKeyPaths(ctx, assetJsDir, JsPath),
	}
}

func newLoaderAssetPath(ctx context.Context, dirSelector assetDirSelector, httpPath func(fileName string) string) loaderAssetPath {
	return func(fileName string, assets *assets.Assets) (string, error) {
		filePath := filepath.Join(dirSelector(assets), fileName)
		file, err := assets.ReadAssetFile(ctx, filePath)
		if err != nil {
			return "", err
		}
		return httpPath(file.FileNameWithHash()), nil
	}
}

func newLoaderAssetPaths(ctx context.Context, dirSelector assetDirSelector, httpPath func(fileName string) string) loaderAssetPaths {
	return func(fileNames []string, assets *assets.Assets) ([]string, error) {
		paths, err := loadFileNamesWithHash(ctx, fileNames, dirSelector, assets)
		if err != nil {
			return nil, err
		}

		for i, fineName := range paths {
			paths[i] = httpPath(fineName)
		}

		return paths, nil
	}
}

func newLoaderAssetNameKeyPaths(
	ctx context.Context,
	dirSelector assetDirSelector,
	httpPath func(fileName string) string,
) loaderAssetNameKeyPaths {
	return func(fileNames []string, assets *assets.Assets) (map[string]string, error) {
		paths, err := loadFileNamesWithHash(ctx, fileNames, dirSelector, assets)
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
