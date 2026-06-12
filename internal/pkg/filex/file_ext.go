package filex

import "path/filepath"

// FileNameWithoutExt returns the file name without its extension
func FileNameWithoutExt(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}
