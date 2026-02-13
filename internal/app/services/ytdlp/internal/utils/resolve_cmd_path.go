package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func ResolveCmdPath(cmdName, binDir string) (string, error) {
	// On Windows, add .exe suffix if missing
	if runtime.GOOS == "windows" && !strings.HasSuffix(cmdName, ".exe") {
		cmdName += ".exe"
	}

	// try PATH
	if path, err := LookupExecutable(cmdName); err == nil {
		return path, nil
	}

	// try config dir
	if binDir != "" {
		cmdPath := filepath.Join(binDir, cmdName)
		if fi, err := os.Stat(cmdPath); err == nil && !fi.IsDir() {
			return cmdPath, nil
		}
	}

	// try executable directory (same folder as your service binary)
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		cmdPath := filepath.Join(exeDir, cmdName)
		if fi, err := os.Stat(cmdPath); err == nil && !fi.IsDir() {
			return cmdPath, nil
		}
	}

	return "", fmt.Errorf(
		"%q executable not found. Tried PATH lookup and config directory %q (full path: %q)",
		cmdName,
		binDir,
		filepath.Join(binDir, cmdName),
	)
}
