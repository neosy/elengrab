package icons

import (
	"context"
	"html/template"
	"os"
	"path/filepath"
)

func FileRaw(fileName string, iconsDir string) template.HTML {
	const svgEmpty = `<svg width="1em" height="1em"></svg>`

	ctx := context.Background()

	icon, _, _ := iconRep.Find(ctx, fileName)
	if icon != nil {
		return icon.raw
	}

	filePath := filepath.Join(iconsDir, fileName)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return svgEmpty
	}

	iconRep.Save(ctx, fileName, &iconEntry{template.HTML(data)})

	return template.HTML(data)
}
