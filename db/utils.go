package database

import (
	"path/filepath"
	"runtime"
)

func filePathForMigrate(p string) string {
	if runtime.GOOS == "windows" {
		return "file:" + p
	}

	p = filepath.ToSlash(p)
	return "file://" + p
}
