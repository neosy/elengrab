package utils

import (
	"os"

	"github.com/neosy/elengrab/internal/pkg/errorx"
)

// CreateTempDir creates a temporary working directory inside baseDir.
// Returns the path to the temp directory and a cleanup function that removes it.
func CreateTempDir(baseDir string, prefix string) (string, func() error, error) {
	// Ensure base directory exists
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return "", nil, errorx.Errorf("failed to create base temp dir: %w", err)
	}

	// Create isolated temp directory for this run
	workDir, err := os.MkdirTemp(baseDir, prefix)
	if err != nil {
		return "", nil, err
	}

	// Cleanup function to remove the temp directory
	cleanup := func() error {
		return os.RemoveAll(workDir)
	}

	return workDir, cleanup, nil
}
