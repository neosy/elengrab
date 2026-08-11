package icons

import (
	"html/template"

	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

var (
	iconsDir           string
	iconFileNamesByKey = make(map[string]string)
)

type Icon struct {
	key      string
	fileName string
}

func newIcon(key, fileName string) Icon {
	iconFileNamesByKey[key] = fileName

	return Icon{
		key:      key,
		fileName: fileName,
	}
}

func (icon Icon) IsZero() bool {
	return icon.key == "" && icon.fileName == ""
}

func (icon Icon) Key() string {
	return icon.key
}

func (icon Icon) FileName() string {
	return icon.fileName
}

func (icon Icon) FileRaw() template.HTML {
	return FileRaw(icon.fileName, iconsDir)
}

func (icon *Icon) URLPath() string {
	return httppaths.BuildIconPath(icon.fileName)
}

func FileNamesByKey() map[string]string {
	return iconFileNamesByKey
}
