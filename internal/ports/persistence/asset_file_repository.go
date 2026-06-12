package persistence

import (
	"context"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
)

type AssetFileCacheRepository interface {
	Save(file *dtypes.AssetFile) error
	SaveNegative(hash, filePath string) error
	Find(hash string) (*dtypes.AssetFile, memsimple.CacheStatus, error)
	FindByFilePath(filePath string) (*dtypes.AssetFile, memsimple.CacheStatus, error)
	Exists(hash string) (bool, error)
	ExistsByFilePath(filePath string) (bool, error)
	CleanExpired(context.Context) error
}
