package assets

import "path/filepath"

var (
	static      = "static"
	staticImg   = filepath.Join(static, "img")
	staticIcons = filepath.Join(staticImg, "icons")
	staticCss   = filepath.Join(static, "css")
	staticJs    = filepath.Join(static, "js")
	staticPwa   = filepath.Join(static, "pwa")
)

type FolderPaths struct {
	assets string
	static string
	img    string
	icons  string
	css    string
	js     string
	pwa    string
}

func NewFolderPaths(assetsDir string) FolderPaths {
	return FolderPaths{
		assets: assetsDir,
		static: filepath.Join(assetsDir, static),
		img:    filepath.Join(assetsDir, staticImg),
		icons:  filepath.Join(assetsDir, staticIcons),
		css:    filepath.Join(assetsDir, staticCss),
		js:     filepath.Join(assetsDir, staticJs),
		pwa:    filepath.Join(assetsDir, staticPwa),
	}
}

func (f *FolderPaths) Assets() string {
	return f.assets
}

func (f *FolderPaths) Static() string {
	return f.static
}

func (f *FolderPaths) Img() string {
	return f.img
}

func (f *FolderPaths) Icons() string {
	return f.icons
}

func (f *FolderPaths) Css() string {
	return f.css
}

func (f *FolderPaths) Js() string {
	return f.js
}

func (f *FolderPaths) Pwa() string {
	return f.pwa
}
