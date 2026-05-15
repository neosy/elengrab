//go:build !windows

package core

import (
	"syscall"

	storagetypes "github.com/neosy/elengrab/internal/infrastructure/storage/filesystem/types"
)

// Stats returns filesystem usage for the given path.
func (s *storage) stats() (storagetypes.StorageStats, error) {
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
