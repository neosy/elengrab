package core

import (
	"syscall"
)

// StorageStats represents disk usage information.
type StorageStats struct {
	Total uint64
	Used  uint64
	Free  uint64
}

// Stats returns filesystem usage for the given path.
func (s *storage) Stats() (StorageStats, error) {
	var stat syscall.Statfs_t

	err := syscall.Statfs(s.basePath, &stat)
	if err != nil {
		return StorageStats{}, err
	}

	// block size
	bsize := uint64(stat.Bsize)

	total := stat.Blocks * bsize
	free := stat.Bfree * bsize
	avail := stat.Bavail * bsize

	// Used is more accurately total - available to non-root
	used := total - avail

	return StorageStats{
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
