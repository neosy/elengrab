//go:build windows
// +build windows

package helper

import (
	"golang.org/x/sys/windows"
)

// DiskUsage for Windows
func DiskUsage(path string) (free uint64, total uint64, used uint64, err error) {
	var freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes uint64

	pathPtr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return
	}

	err = windows.GetDiskFreeSpaceEx(pathPtr, &freeBytesAvailable, &totalNumberOfBytes, &totalNumberOfFreeBytes)
	if err != nil {
		return
	}

	free = freeBytesAvailable
	total = totalNumberOfBytes
	used = total - free
	return
}
