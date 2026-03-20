//go:build linux || darwin
// +build linux darwin

package helper

import "syscall"

// DiskUsage for Unix systems
func DiskUsage(path string) (free uint64, total uint64, used uint64, err error) {
	var stat syscall.Statfs_t
	err = syscall.Statfs(path, &stat)
	if err != nil {
		return
	}
	total = stat.Blocks * uint64(stat.Bsize)
	free = stat.Bfree * uint64(stat.Bsize)
	used = total - free
	return
}
