package iconfig

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// defaultRootDir returns the default application data directory for the current OS.
func defaultRootDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if runtime.GOOS == "windows" {
		// Use %LOCALAPPDATA% when available.
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData != "" {
			return filepath.Join(localAppData, AppName), nil
		}

		// Fall back to the user's local application data directory.
		return filepath.Join(home, "AppData", "Local", AppName), nil
	}

	// Use a hidden directory in the user's home directory on Linux and macOS.
	return filepath.Join(home, "."+strings.ToLower(AppName)), nil
}
