package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func ResolveCmdPath(cmdName, binDir string) (string, error) {
	cmdName = normalizeCommandName(cmdName)

	// Try configured directory
	if binDir != "" {
		cmdPath := filepath.Join(binDir, cmdName)
		if fi, err := os.Stat(cmdPath); err == nil && !fi.IsDir() {
			return cmdPath, nil
		}
	}

	// Try executable directory (same folder as the service binary)
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		cmdPath := filepath.Join(exeDir, cmdName)

		if fi, err := os.Stat(cmdPath); err == nil && !fi.IsDir() {
			return cmdPath, nil
		}
	}

	// Try PATH
	if path, err := LookupExecutable(cmdName); err == nil {
		return path, nil
	}

	return "", fmt.Errorf("%s executable not found", cmdName)
}
