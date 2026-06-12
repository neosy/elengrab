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
	cacheByHash     memsimple.Cache[string, dtypes.AssetFile]
	cacheByFilePath memsimple.Cache[string, dtypes.AssetFile]
}

// newAssetFileRepository returns a new object for the repository
func newAssetFileRepository(ttl time.Duration) *AssetFileRepository {
	r := &AssetFileRepository{
		cacheByHash:     memsimple.NewCacheWithDeaultCopier[string, dtypes.AssetFile, *dtypes.AssetFile](),
		cacheByFilePath: memsimple.NewCacheWithDeaultCopier[string, dtypes.AssetFile, *dtypes.AssetFile](),
	}
	r.Repository.Init(ttl)
	return r
}

func (r *AssetFileRepository) Save(assetFile *dtypes.AssetFile) error {
	if assetFile == nil {
		return nil
	}

	save := func() error {
		r.cacheByHash.Save(
			assetFile.Hash,
			assetFile,
			r.TTL(),
		)
		r.cacheByFilePath.Save(
			assetFile.FilePath,
			assetFile,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(save)
}

func (r *AssetFileRepository) SaveNegative(hash, filePath string) error {
	if hash == "" {
		return nil
	}

	save := func() error {
		r.cacheByHash.Save(
			hash,
			nil,
			r.TTL(),
		)
		r.cacheByFilePath.Save(
			filePath,
			nil,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(save)
}

// Delete removes a asset file from the in-memory repository using its ID.
func (r *AssetFileRepository) Delete(hash string) error {
	delete := func() error {
		if hash != "" {
			file := r.cacheByHash.Find(hash)
			r.cacheByHash.Delete(hash)
			r.cacheByFilePath.Delete(file.FilePath)
		}
		return nil
	}
	return r.Repository.Delete(delete)
}

// Find retrieves a file by its hash from the repository.
func (r *AssetFileRepository) Find(hash string) (*dtypes.AssetFile, memsimple.CacheStatus, error) {
	find := func() (*dtypes.AssetFile, memsimple.CacheStatus, error) {
		data, status := r.cacheByHash.FindWithStatus(hash)
		return data, status, nil
	}
	return r.Repository.FindWithStatus(find)
}

// FindByFilePath retrieves a file by its file path from the repository.
func (r *AssetFileRepository) FindByFilePath(filePath string) (*dtypes.AssetFile, memsimple.CacheStatus, error) {
	find := func() (*dtypes.AssetFile, memsimple.CacheStatus, error) {
		data, status := r.cacheByFilePath.FindWithStatus(filePath)
		return data, status, nil
	}
	return r.Repository.FindWithStatus(find)
}

// Exists checks if a asset file exists by its hash.
func (r *AssetFileRepository) Exists(hash string) (bool, error) {
	exists := func() (bool, error) {
		return r.cacheByHash.Exists(hash), nil
	}
	return r.Repository.Exists(exists)
}

// ExistsByFilePath checks if a asset file exists by its file path.
func (r *AssetFileRepository) ExistsByFilePath(filePath string) (bool, error) {
	exists := func() (bool, error) {
		return r.cacheByFilePath.Exists(filePath), nil
	}
	return r.Repository.Exists(exists)
}

// CleanExpired cleans expired entries from the repository.
func (r *AssetFileRepository) CleanExpired(context.Context) error {
	// Define a clean function to remove expired entries from the cache.
	clean := func() error {
		r.cacheByHash.CleanExpired()
		r.cacheByFilePath.CleanExpired()
		return nil
	}
	// Call the base repository's CleanExpired method with the custom clean function.
	return r.Repository.CleanExpired(clean)
}
