package core

import (
	"fmt"
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
