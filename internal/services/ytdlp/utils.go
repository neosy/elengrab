package ytdlpsrv

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func resolveCmdPath(cmdName, binDir string) (string, error) {
	if path, err := exec.LookPath(cmdName); err == nil {
		return path, nil
	}

	cmdPath := filepath.Join(binDir, cmdName)
	if fi, err := os.Stat(cmdPath); err == nil && !fi.IsDir() {
		return cmdPath, nil
	}

	return "", fmt.Errorf("%s not found: tried config path %q and PATH lookup", cmdName, cmdPath)
}

func checkDir(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %v", err)
		} else {
			return err
		}
	} else if !info.IsDir() {
		return errors.New("path exists but is not a directory")
	}

	return nil
}

func fileNameWithoutExt(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}
