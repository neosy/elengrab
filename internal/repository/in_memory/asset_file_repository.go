package inmemory

import (
	"context"
	"time"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
)

// Defines the structure for the in-memory repository.
type AssetFileRepository struct {
	// Embeds the base Repository
	memsimple.Repository[dtypes.AssetFile]

	// Cache for storing Files by their unique Hash.
	cacheByPath memsimple.Cache[string, dtypes.AssetFile]
}

// newAssetFileRepository returns a new object for the repository
func newAssetFileRepository(ttl time.Duration) *AssetFileRepository {
	r := &AssetFileRepository{
		cacheByPath: memsimple.NewCacheWithDeaultCopier[string, dtypes.AssetFile, *dtypes.AssetFile](),
	}
	r.Repository.Init(ttl)
	return r
}

func (r *AssetFileRepository) Save(assetFile *dtypes.AssetFile) error {
	if assetFile == nil {
		return nil
	}

	save := func() error {
		r.cacheByPath.Save(
			assetFile.FilePath,
			assetFile,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(save)
}

func (r *AssetFileRepository) SaveNegative(filePath string) error {
	if filePath == "" {
		return nil
	}

	save := func() error {
		r.cacheByPath.Save(
			filePath,
			nil,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(save)
}

// Delete removes a asset file from the in-memory repository using its filePath.
func (r *AssetFileRepository) Delete(filePath string) error {
	delete := func() error {
		if filePath != "" {
			r.cacheByPath.Delete(filePath)
		}
		return nil
	}
	return r.Repository.Delete(delete)
}

// FindByPath retrieves a file by its file path from the repository.
func (r *AssetFileRepository) FindByPath(filePath string) (*dtypes.AssetFile, memsimple.CacheStatus, error) {
	find := func() (*dtypes.AssetFile, memsimple.CacheStatus, error) {
		data, status := r.cacheByPath.FindWithStatus(filePath)
		return data, status, nil
	}
	return r.Repository.FindWithStatus(find)
}

// ExistsByPath checks if a asset file exists by its file path.
func (r *AssetFileRepository) ExistsByPath(filePath string) (bool, error) {
	exists := func() (bool, error) {
		return r.cacheByPath.Exists(filePath), nil
	}
	return r.Repository.Exists(exists)
}

// CleanExpired cleans expired entries from the repository.
func (r *AssetFileRepository) CleanExpired(context.Context) error {
	// Define a clean function to remove expired entries from the cache.
	clean := func() error {
		r.cacheByPath.CleanExpired()
		return nil
	}
	// Call the base repository's CleanExpired method with the custom clean function.
	return r.Repository.CleanExpired(clean)
}
