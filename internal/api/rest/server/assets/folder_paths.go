package assets

import (
	"path/filepath"
)

var (
	static       = "static"
	staticCss    = filepath.Join(static, "css")
	staticFonts  = filepath.Join(static, "fonts")
	staticJs     = filepath.Join(static, "js")
	staticImages = filepath.Join(static, "images")
	staticIcons  = filepath.Join(static, "icons")
	staticPwa    = filepath.Join(static, "pwa")

	templates = "templates"
	layouts   = filepath.Join(templates, "layouts")
	pages     = filepath.Join(templates, "pages")
)

func NewFolderPaths(assetsDir string) FolderPaths {
	return FolderPaths{
		Assets: makeFolderPath(assetsDir),

		Static: makeFolderPath(filepath.Join(assetsDir, static)),
		Css:    makeFolderPath(filepath.Join(assetsDir, staticCss)),
		Fonts:  makeFolderPath(filepath.Join(assetsDir, staticFonts)),
		Js:     makeFolderPath(filepath.Join(assetsDir, staticJs)),
		Img:    makeFolderPath(filepath.Join(assetsDir, staticImages)),
		Icons:  makeFolderPath(filepath.Join(assetsDir, staticIcons)),
		Pwa:    makeFolderPath(filepath.Join(assetsDir, staticPwa)),

		Templates: makeFolderPath(filepath.Join(assetsDir, templates)),
		Layouts:   makeFolderPath(filepath.Join(assetsDir, layouts)),
		Pages:     makeFolderPath(filepath.Join(assetsDir, pages)),
	}
}

type (
	folderPath func() string

	FolderPaths struct {
		Assets folderPath

		Static folderPath
		Css    folderPath
		Fonts  folderPath
		Js     folderPath
		Img    folderPath
		Icons  folderPath
		Pwa    folderPath

		Templates folderPath
		Layouts   folderPath
		Pages     folderPath
	}
)

func makeFolderPath(dir string) folderPath {
	return func() string {
		return dir
	}
}
