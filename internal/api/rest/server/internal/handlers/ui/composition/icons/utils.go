package icons

import (
	"html/template"
	"os"
	"path/filepath"
)

func FileRaw(fileName string, iconsDir string) template.HTML {
	const svgEmpty = `<svg width="1em" height="1em"></svg>`

	icon := iconCache.Find(fileName)
	if icon != nil {
		return icon.raw
	}

	filePath := filepath.Join(iconsDir, fileName)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return svgEmpty
	}

	iconCache.Save(fileName, &iconEntry{template.HTML(data)}, iconCacheTTL)

	return template.HTML(data)
}
