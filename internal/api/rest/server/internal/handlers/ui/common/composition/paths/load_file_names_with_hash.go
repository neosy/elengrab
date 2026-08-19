package paths

import (
	"context"
	"path/filepath"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
)

func loadFileNamesWithHash(
	ctx context.Context,
	fileNames []string,
	dirSelector assetDirSelector,
	assets *assets.Assets,
) ([]string, error) {
	var result = make([]string, len(fileNames))

	for i, f := range fileNames {
		var err error
		file, err := assets.ReadAssetFile(ctx, filepath.Join(dirSelector(assets), f))
		if err != nil {
			return nil, err
		}
		result[i] = file.FileNameWithHash()
	}

	return result, nil
}
