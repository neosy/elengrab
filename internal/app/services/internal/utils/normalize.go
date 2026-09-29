package utils

import (
	"runtime"
	"strings"
)

func normalizeCommandName(cmdName string) string {
	// On Windows, add .exe suffix if missing
	if runtime.GOOS == "windows" && !strings.HasSuffix(cmdName, ".exe") {
		return cmdName + ".exe"
	}

	return cmdName
}
