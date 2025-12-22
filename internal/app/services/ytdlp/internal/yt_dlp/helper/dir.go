package helper

import (
	"fmt"
	"os"
)

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
