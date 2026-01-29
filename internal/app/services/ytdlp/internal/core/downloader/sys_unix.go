//go:build !windows
// +build !windows

package downloader

import (
	"os/exec"
	"syscall"
)

func newSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setpgid: true,
	}
}

func tryGracefulKill(pgid int) {
	_ = syscall.Kill(pgid, syscall.SIGTERM)
}

func forceKill(cmd *exec.Cmd) {
	if cmd.Process != nil {
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
}
