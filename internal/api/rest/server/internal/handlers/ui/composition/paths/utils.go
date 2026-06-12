package paths

import (
	"path/filepath"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
)

func fileNamesWithHash(dir string, fileNames []string) ([]string, error) {
	var result = make([]string, len(fileNames))

	for i, f := range fileNames {
		var err error
		file, err := assets.ReadAssetFileWithHash(filepath.Join(dir, f))
		if err != nil {
			return nil, err
		}
		result[i] = file.FileNameWithHash()
	}

	return result, nil
}
