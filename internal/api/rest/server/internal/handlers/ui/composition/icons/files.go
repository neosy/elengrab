package icons

import (
	"html/template"
	"os"
	"path/filepath"
)

func FileNames() map[string]any {
	return iconFileNames
}

func FileName(key string) string {
	return iconFileNames[key].(string)
}

func FileRaw(fileName string, svgDir string) template.HTML {
	const svgEmpty = `<svg width="1em" height="1em"></svg>`

	icon := iconCache.Find(fileName)
	if icon != nil {
		return icon.raw
	}

	filePath := filepath.Join(svgDir, fileName)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return svgEmpty
	}

	iconCache.Save(fileName, &iconEntry{template.HTML(data)}, iconCacheTTL)

	return template.HTML(data)
}

func FileRawByKey(key string, svgDir string) template.HTML {
	return FileRaw(FileName(key), svgDir)
}
