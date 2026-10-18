package pstorage

import storagetypes "github.com/neosy/elengrab/internal/infrastructure/storage/filesystem/types"

// DownloadsStorage defines the interface for managing downloads storage.
type DownloadsStorage interface {
	// Exists checks if file exists in storage.
	Exists(uniqueFileName string) (bool, error)

	// Move moves the file at fromFilePath into storage under uniqueFileName.
	// The source file is removed on success.
	// Returns an error if a file with this name already exists.
	Move(fromFilePath, uniqueFileName string) error

	// Delete removes the file stored under uniqueFileName.
	Delete(uniqueFileName string) error

	// BasePath returns the root directory of the downloads storage.
	BasePath() string

	// MediaPath returns the directory where downloaded media files are stored.
	MediaPath() string

	// Path returns the filesystem path where the file with uniqueFileName is (or would be) stored.
	// It does not check that the file exists.
	Path(uniqueFileName string) string

	// Stats returns filesystem usage for the given path.
	Stats() (storagetypes.StorageStats, error)

	// Used returns the total disk space consumed by all files under the base path.
	// It recursively walks the directory and sums file sizes.
	Used() (uint64, error)
}
