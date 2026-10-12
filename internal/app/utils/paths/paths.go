package upaths

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	nfile "github.com/neosy/elengrab/internal/pkg/filex"
)

// absPath resolves a relative path to an absolute path using current working directory.
// Exits the program if resolving fails.
func AbsPath(root, path string) string {
	path, err := nfile.AbsPath(root, path)
	if err != nil {
		log.Fatal(err.Error())
	}
	return path
}

// ExecutableDir returns the directory containing the running executable.
func ExecutableDir() (string, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("get executable path: %w", err)
	}

	return filepath.Dir(executablePath), nil
}

// EnsureDir checks whether the directory exists and creates it if necessary.
func EnsureDir(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("cannot access directory %s: %w", path, err)
		}

		err := os.MkdirAll(path, 0o755)
		if err != nil {
			return fmt.Errorf("cannot create directory %s: %w", path, err)
		}

		return nil
	}

	return nfile.CheckDir(path)
}
