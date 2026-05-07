package pstorage

import storagetypes "github.com/neosy/elengrab/internal/infrastructure/storage/filesystem/types"

// DownloadsStorage defines the interface for managing downloads storage.
type DownloadsStorage interface {
	// Exists checks if file exists in storage.
	Exists(uniqueFileName string) (bool, error)

	// Move file from directory to storage
	Move(fromFilePath, uniqueFileName string) error

	// Delete removes the dowload data.
	Delete(uniqueFileName string) error

	// BasePath
	BasePath() string

	// MediaPath
	MediaPath() string

	// Path returns the full file path.
	Path(uniqueFileName string) string

	// Stats returns filesystem usage for the given path.
	Stats() (storagetypes.StorageStats, error)

	// Used returns the total disk space consumed by all files under the base path.
	// It recursively walks the directory and sums file sizes.
	Used() (uint64, error)
}
