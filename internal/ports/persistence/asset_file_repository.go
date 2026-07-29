package persistence

import (
	"context"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/cache/memory"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
)

type AssetFileCacheRepository interface {
	memory.CacheRepository

	Save(file *dtypes.AssetFile) error
	SaveNegative(filePath string) error

	FindByPath(filePath string) (*dtypes.AssetFile, memsimple.CacheStatus, error)
	ExistsByPath(filePath string) (bool, error)

	CleanExpired(context.Context) error
}
