package core

import "syscall"

// StorageStats represents disk usage information.
type StorageStats struct {
	Total int64
	Used  int64
	Free  int64
}

// Stats returns filesystem usage for the given path.
func (s *storage) Stats() (StorageStats, error) {
	var stat syscall.Statfs_t

	err := syscall.Statfs(s.basePath, &stat)
	if err != nil {
		return StorageStats{}, err
	}

	// block size
	bsize := int64(stat.Bsize)

	total := int64(stat.Blocks) * bsize
	free := int64(stat.Bfree) * bsize
	avail := int64(stat.Bavail) * bsize

	// Used is more accurately total - available to non-root
	used := total - avail

	return StorageStats{
		Total: total,
		Used:  used,
		Free:  free,
	}, nil
}
