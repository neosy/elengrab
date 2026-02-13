package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func CheckFFmpeg(ffmpegName string) error {
	if runtime.GOOS == "windows" && !strings.HasSuffix(ffmpegName, ".exe") {
		ffmpegName += ".exe"
	}

	// Ensure ffmpeg is available in PATH
	cmd := exec.Command(ffmpegName, "-version")
	if err := cmd.Run(); err == nil {
		return nil
	}

	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		cmdPath := filepath.Join(exeDir, ffmpegName)

		cmd = exec.Command(cmdPath, "-version")
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	return fmt.Errorf("%s not found in PATH. Please install %s and add it to PATH", ffmpegName, ffmpegName)
}
