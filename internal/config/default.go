package iconfig

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func defaultRootDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if runtime.GOOS == "windows" {
		// %LOCALAPPDATA%
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData != "" {
			return filepath.Join(localAppData, AppName), nil
		}
		// fallback
		return filepath.Join(home, "AppData", "Local", AppName), nil
	}

	// Linux / macOS
	return filepath.Join(home, "."+strings.ToLower(AppName)), nil
}
