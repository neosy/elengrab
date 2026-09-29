package utils

import (
	"os/exec"
)

func LookupExecutable(cmdName string) (string, error) {
	return exec.LookPath(normalizeCommandName(cmdName))
}
