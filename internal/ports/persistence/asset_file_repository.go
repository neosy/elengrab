package persistence

import (
	"context"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/cache/memory"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
)

type AssetFileCacheRepositoryFactory func() AssetFileCacheRepository

type AssetFileCacheRepository interface {
	memory.CacheRepository

	Save(ctx context.Context, file *dtypes.AssetFile) error
	SaveNegative(ctx context.Context, filePath string) error

	FindByPath(ctx context.Context, filePath string) (*dtypes.AssetFile, memsimple.CacheStatus, error)
	ExistsByPath(ctx context.Context, filePath string) (bool, error)

	CleanExpired(context.Context) error
}
