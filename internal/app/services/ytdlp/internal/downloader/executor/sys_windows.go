//go:build windows
// +build windows

package executor

import (
	"os/exec"
	"syscall"
)

func newSysProcAttr() *syscall.SysProcAttr {
	// Windows does not support Setpgid.
	return &syscall.SysProcAttr{}
}

func tryGracefulKill(pgid int) {
	// Windows does not support process group kill via syscall.
	// No-op on Windows.
}

func forceKill(cmd *exec.Cmd) {
	if cmd.Process != nil {
		_ = cmd.Process.Kill()
	}
}
