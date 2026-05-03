package pstorage

import fsstorage "github.com/neosy/elengrab/internal/infrastructure/storage/filesystem"

// DownloadsStorage defines the interface for managing downloads storage.
type DownloadsStorage interface {
	// Exists checks if file exists in storage.
	Exists(uniqueFileName string) (bool, error)

	// Move file from directory to storage
	Move(fromFilePath, uniqueFileName string) error

	// Delete removes the dowload data.
	Delete(uniqueFileName string) error

	// Path returns the full file path.
	Path(uniqueFileName string) string

	// BasePath
	BasePath() string

	// Stats returns filesystem usage for the given path.
	Stats() (fsstorage.StorageStats, error)
}
