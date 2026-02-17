//go:build !windows
// +build !windows

package executor

import (
	"errors"
	"fmt"
	"os/exec"
	"syscall"
)

func newSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setpgid: true,
	}
}

func tryGracefulKill(pgid int) error {
	err := syscall.Kill(pgid, syscall.SIGTERM)
	if err != nil {
		if !errors.Is(err, syscall.ESRCH) {
			return fmt.Errorf("failed SIGTERM pgid %v: %w", pgid, err)
		}
	}
	return nil
}

func forceKill(cmd *exec.Cmd) error {
	if cmd.Process != nil {
		return nil
	}

	pgid := cmd.Process.Pid
	err := syscall.Kill(-pgid, syscall.SIGKILL)
	if err != nil {
		if !errors.Is(err, syscall.ESRCH) {
			return fmt.Errorf("failed SIGKILL pgid %v: %w", pgid, err)
		}
	}

	return nil
}
