package nfile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// AbsPathCwd returns the absolute path for the given path.
// Relative paths are resolved relative to the current working directory.
// Example:
//
//	path, err := AbsPathCwd("./downloads")
//	// path -> "/home/user/project/downloads"
func AbsPathCwd(path string) (string, error) {
	if path == "" {
		return "", errors.New("path is empty")
	}

	if filepath.IsAbs(path) {
		return path, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot get working directory: %v", err)
	}
	return filepath.Join(cwd, path), nil
}

// AbsPath returns the absolute path for the given path.
// If the path is relative, it is joined with the provided root directory.
// If the path is already absolute, it is returned as is.
// Example:
//
//	appRoot := "/app_n"
//	path, err := AbsPath(appRoot, "./downloads")
//	// path -> "/app_n/downloads"
func AbsPath(root, path string) (string, error) {
	if path == "" {
		return "", errors.New("path is empty")
	}

	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}

	if root == "" {
		return "", errors.New("root path is empty")
	}

	if !filepath.IsAbs(root) {
		return "", fmt.Errorf("root must be absolute: %s", root)
	}

	if filepath.IsAbs(path) {
		return path, nil
	}

	return filepath.Clean(filepath.Join(root, path)), nil
}
