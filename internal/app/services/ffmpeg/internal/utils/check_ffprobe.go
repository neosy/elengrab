package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func CheckFFprobe(ffprobeName string) error {
	if runtime.GOOS == "windows" && !strings.HasSuffix(ffprobeName, ".exe") {
		ffprobeName += ".exe"
	}

	// Ensure ffprobe is available in PATH
	cmd := exec.Command(ffprobeName, "-version")
	if err := cmd.Run(); err == nil {
		return nil
	}

	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		cmdPath := filepath.Join(exeDir, ffprobeName)

		cmd = exec.Command(cmdPath, "-version")
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	return fmt.Errorf("%s not found in PATH. Please install %s and add it to PATH", ffprobeName, ffprobeName)
}
