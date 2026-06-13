package assets

import "github.com/neosy/elengrab/internal/ports/persistence"

type Assets struct {
	folderPaths FolderPaths
	fileCache   persistence.AssetFileCacheRepository

	// settings
	assetsDirPath string
}

func NewAssets(
	assetsDirPath string,
	fileCache persistence.AssetFileCacheRepository,
) *Assets {
	return &Assets{
		folderPaths: newFolderPaths(assetsDirPath),
		fileCache:   fileCache,

		// settings
		assetsDirPath: assetsDirPath,
	}
}

func (a *Assets) FolderPaths() FolderPaths {
	return a.folderPaths
}

func (a *Assets) DirPath() string {
	return a.assetsDirPath
}
