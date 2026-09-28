package utils

import (
	"os/exec"
	"runtime"
	"strings"
)

func LookupExecutable(cmdName string) (string, error) {
	if runtime.GOOS == "windows" && !strings.HasSuffix(cmdName, ".exe") {
		cmdName += ".exe"
	}
	return exec.LookPath(cmdName)
}
