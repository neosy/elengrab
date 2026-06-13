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

func newFolderPaths(assetsDirPath string) FolderPaths {
	return FolderPaths{
		Assets: makeFolderPath(assetsDirPath),

		Static: makeFolderPath(filepath.Join(assetsDirPath, static)),
		Css:    makeFolderPath(filepath.Join(assetsDirPath, staticCss)),
		Fonts:  makeFolderPath(filepath.Join(assetsDirPath, staticFonts)),
		Js:     makeFolderPath(filepath.Join(assetsDirPath, staticJs)),
		Img:    makeFolderPath(filepath.Join(assetsDirPath, staticImages)),
		Icons:  makeFolderPath(filepath.Join(assetsDirPath, staticIcons)),
		Pwa:    makeFolderPath(filepath.Join(assetsDirPath, staticPwa)),

		Templates: makeFolderPath(filepath.Join(assetsDirPath, templates)),
		Layouts:   makeFolderPath(filepath.Join(assetsDirPath, layouts)),
		Pages:     makeFolderPath(filepath.Join(assetsDirPath, pages)),
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
