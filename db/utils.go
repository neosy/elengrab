package database

import (
	"path/filepath"
	"runtime"
	"strings"
)

func filePathForMigrate(p string) string {
	if runtime.GOOS == "windows" && !strings.HasPrefix(p, "/") {
		return p
	}
	p = filepath.ToSlash(p)
	return "file://" + p
}
