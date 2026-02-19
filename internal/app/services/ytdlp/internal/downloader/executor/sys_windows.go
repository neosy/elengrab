//go:build windows

package executor

import (
	"os/exec"
	"syscall"
)

func newSysProcAttr() *syscall.SysProcAttr {
	// Windows does not support Setpgid.
	return &syscall.SysProcAttr{}
}

func tryGracefulKill(int) error {
	// Windows does not support process group kill via syscall.
	// No-op on Windows.
	return nil
}

func forceKill(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return nil
	}
	return cmd.Process.Kill()
}
