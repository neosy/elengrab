//go:build windows

package core

import (
	"golang.org/x/sys/windows"

	storagetypes "github.com/neosy/elengrab/internal/infrastructure/storage/filesystem/types"
)

func (s *storage) stats() (storagetypes.StorageStats, error) {
	pathPtr, err := windows.UTF16PtrFromString(s.basePath)
	if err != nil {
		return storagetypes.StorageStats{}, err
	}

	var freeBytesAvailable uint64
	var totalNumberOfBytes uint64
	var totalNumberOfFreeBytes uint64

	err = windows.GetDiskFreeSpaceEx(
		pathPtr,
		&freeBytesAvailable,
		&totalNumberOfBytes,
		&totalNumberOfFreeBytes,
	)
	if err != nil {
		return storagetypes.StorageStats{}, err
	}

	used := totalNumberOfBytes - freeBytesAvailable

	return storagetypes.StorageStats{
		Total: totalNumberOfBytes,
		Used:  used,
		Free:  totalNumberOfFreeBytes,
	}, nil
}
