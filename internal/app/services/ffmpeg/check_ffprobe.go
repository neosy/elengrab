package ffmpegsrv

import (
	"fmt"
	"os"
	"os/exec"
)

func checkFFprobe(cmdPath string) error {
	fi, err := os.Stat(cmdPath)
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return fmt.Errorf("%q is a directory", cmdPath)
	}

	// Ensure ffprobe is available in PATH
	cmd := exec.Command(cmdPath, "-version")
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func checkDeno(cmdPath string) error {
	fi, err := os.Stat(cmdPath)
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return fmt.Errorf("%q is a directory", cmdPath)
	}

	cmd := exec.Command(cmdPath, "--version")
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
