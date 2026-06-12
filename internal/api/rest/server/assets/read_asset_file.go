package assets

import (
	"os"
	"path/filepath"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/assetx"
)

func readOriginalFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func readOriginalFileWithHash(filePath string) (*dtypes.AssetFile, error) {
	data, err := readOriginalFile(filePath)
	if err != nil {
		return nil, err
	}

	hash := HashFromData(data)

	return dtypes.NewAssetFile(filePath, hash, data), nil
}

func ReadJsFileWithHashedImports(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	updatedData := assetx.RewriteJSImports(
		data,
		func(name string) string {
			filePath := filepath.Join(filepath.Dir(filePath), name)
			file, err := readOriginalFileWithHash(filePath)
			if err != nil {
				return name
			}
			return file.FileNameWithHash()
		},
	)

	return updatedData, nil
}

func ReadJsFileWithHashedImportsAndHash(filePath string) (*dtypes.AssetFile, error) {
	data, err := ReadJsFileWithHashedImports(filePath)
	if err != nil {
		return nil, err
	}

	hash := HashFromData(data)

	return dtypes.NewAssetFile(filePath, hash, data), nil
}

func ReadAssetFile(filePath string) ([]byte, error) {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".js":
		return ReadJsFileWithHashedImports(filePath)
	default:
		return readOriginalFile(filePath)
	}
}

func ReadAssetFileWithHash(filePath string) (*dtypes.AssetFile, error) {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".js":
		return ReadJsFileWithHashedImportsAndHash(filePath)
	default:
		return readOriginalFileWithHash(filePath)
	}
}
