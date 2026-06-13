package assets

import (
	"errors"
	"os"
	"path/filepath"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/assetx"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (a *Assets) readOriginalFileNoCache(filePath string) (*dtypes.AssetFile, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	hash := HashFromData(data)

	return dtypes.NewAssetFile(filePath, hash, data), nil
}

func (a *Assets) readOriginalFile(filePath string) (*dtypes.AssetFile, error) {
	file, status, err := a.fileCache.FindByPath(filePath)
	if err != nil {
		return nil, err
	}

	if status == memsimple.CacheStatusNegativeHit {
		return nil, errorx.New("file not found", exceptionx.NOT_FOUND)
	}

	if file != nil {
		return file, nil
	}

	file, err = a.readOriginalFileNoCache(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			a.fileCache.SaveNegative(filePath)
			return nil, errorx.New("file not found", exceptionx.NOT_FOUND)
		}
		return nil, err
	}

	a.fileCache.Save(file)

	return file, nil
}

func (a *Assets) readJsFileHashedImportsNoCache(filePath string) (*dtypes.AssetFile, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	updatedData := assetx.RewriteJSImports(
		data,
		func(name string) string {
			filePath := filepath.Join(filepath.Dir(filePath), name)
			file, err := a.readOriginalFileNoCache(filePath)
			if err != nil {
				return name
			}
			return file.FileNameWithHash()
		},
	)

	hash := HashFromData(updatedData)

	return dtypes.NewAssetFile(filePath, hash, updatedData), nil
}

func (a *Assets) readJsFileHashedImports(filePath string) (*dtypes.AssetFile, error) {
	file, status, err := a.fileCache.FindByPath(filePath)
	if err != nil {
		return nil, err
	}

	if status == memsimple.CacheStatusNegativeHit {
		return nil, errorx.New("file not found", exceptionx.NOT_FOUND)
	}

	if file != nil {
		return file, nil
	}

	file, err = a.readJsFileHashedImportsNoCache(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			a.fileCache.SaveNegative(filePath)
			return nil, errorx.New("file not found", exceptionx.NOT_FOUND)
		}
		return nil, err
	}

	a.fileCache.Save(file)

	return file, nil
}

func (a *Assets) ReadAssetFileNoCache(filePath string) (*dtypes.AssetFile, error) {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".js":
		return a.readJsFileHashedImportsNoCache(filePath)
	default:
		return a.readOriginalFileNoCache(filePath)
	}
}

func (a *Assets) ReadAssetFile(filePath string) (*dtypes.AssetFile, error) {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".js":
		return a.readJsFileHashedImports(filePath)
	default:
		return a.readOriginalFile(filePath)
	}
}
