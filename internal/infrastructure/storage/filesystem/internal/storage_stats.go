package core

import storagetypes "github.com/neosy/elengrab/internal/infrastructure/storage/filesystem/types"

// Stats returns filesystem usage for the given path.
func (s *storage) Stats() (storagetypes.StorageStats, error) {
	return s.stats()
}

// Used returns the total disk space consumed by all files under the base path.
// It recursively walks the directory and sums file sizes.
func (s *storage) Used() (uint64, error) {
	return folderSize(s.BasePath())
}
