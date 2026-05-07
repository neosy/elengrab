package core

import (
	"syscall"

	storagetypes "github.com/neosy/elengrab/internal/infrastructure/storage/filesystem/types"
)

// Stats returns filesystem usage for the given path.
func (s *storage) Stats() (storagetypes.StorageStats, error) {
	var stat syscall.Statfs_t

	err := syscall.Statfs(s.basePath, &stat)
	if err != nil {
		return storagetypes.StorageStats{}, err
	}

	// block size
	bsize := uint64(stat.Bsize)

	total := stat.Blocks * bsize
	free := stat.Bfree * bsize
	avail := stat.Bavail * bsize

	// Used is more accurately total - available to non-root
	used := total - avail

	return storagetypes.StorageStats{
		Total: total,
		Used:  used,
		Free:  free,
	}, nil
}

// Used returns the total disk space consumed by all files under the base path.
// It recursively walks the directory and sums file sizes.
func (s *storage) Used() (uint64, error) {
	return folderSize(s.BasePath())
}
