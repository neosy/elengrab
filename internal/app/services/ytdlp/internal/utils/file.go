package iutils

import "path/filepath"

func FileNameWithoutExt(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}
