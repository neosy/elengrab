package core

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	nfile "github.com/neosy/elengrab/internal/pkg/filex"
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

	err := filepath.WalkDir(path, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		size += uint64(info.Size())
		return nil
	})

	if err != nil {
		return 0, err
	}

	return size, nil
}
