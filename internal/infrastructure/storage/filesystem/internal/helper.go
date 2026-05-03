package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	nfile "github.com/neosy/elengrab/internal/pkg/file"
)

func validateDirPath(path string) error {
	if strings.HasSuffix(path, "/") || strings.HasSuffix(path, "\\") {
		return fmt.Errorf("directory must not end with a slash or backslash: %s", path)
	}

	if err := nfile.CheckDir(path); err != nil {
		return err
	}

	return nil
}

func folderSize(path string) (uint64, error) {
	var size uint64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += uint64(info.Size())
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return size, nil
}
