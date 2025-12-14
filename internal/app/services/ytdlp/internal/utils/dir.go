package iutils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func ResolveCmdPath(cmdName, binDir string) (string, error) {
	if path, err := exec.LookPath(cmdName); err == nil {
		return path, nil
	}

	cmdPath := filepath.Join(binDir, cmdName)
	if fi, err := os.Stat(cmdPath); err == nil && !fi.IsDir() {
		return cmdPath, nil
	}

	return "", fmt.Errorf("%s not found: tried config path %q and PATH lookup", cmdName, cmdPath)
}

func CheckDir(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %v", err)
		} else {
			return err
		}
	} else if !info.IsDir() {
		return errors.New("path exists but is not a directory")
	}

	return nil
}

// CreateTempDir creates a temporary working directory inside baseDir.
// Returns the path to the temp directory and a cleanup function that removes it.
func CreateTempDir(baseDir string, prefix string) (string, func(), error) {
	// Ensure base directory exists
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return "", nil, fmt.Errorf("failed to create base temp dir: %w", err)
	}

	// Create isolated temp directory for this run
	workDir, err := os.MkdirTemp(baseDir, prefix)
	if err != nil {
		return "", nil, err
	}

	// Cleanup function to remove the temp directory
	cleanup := func() {
		_ = os.RemoveAll(workDir)
	}

	return workDir, cleanup, nil
}
