package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func CheckDeno(denoName string) error {
	if runtime.GOOS == "windows" && !strings.HasSuffix(denoName, ".exe") {
		denoName += ".exe"
	}

	// Ensure deno is available in PATH
	cmd := exec.Command(denoName, "--version")
	if err := cmd.Run(); err == nil {
		return nil
	}

	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		cmdPath := filepath.Join(exeDir, denoName)

		cmd = exec.Command(cmdPath, "--version")
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	return fmt.Errorf("%s not found in PATH. Please install %s and add it to PATH", denoName, denoName)
}
