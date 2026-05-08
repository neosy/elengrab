package assets

import "path/filepath"

var (
	static      = "static"
	staticCss   = filepath.Join(static, "css")
	staticFonts = filepath.Join(static, "fonts")
	staticJs    = filepath.Join(static, "js")
	staticImg   = filepath.Join(static, "img")
	staticIcons = filepath.Join(staticImg, "icons")
	staticPwa   = filepath.Join(static, "pwa")
)

type FolderPaths struct {
	assets string
	static string
	css    string
	fonts  string
	js     string
	img    string
	icons  string
	pwa    string
}

func NewFolderPaths(assetsDir string) FolderPaths {
	return FolderPaths{
		assets: assetsDir,
		static: filepath.Join(assetsDir, static),
		css:    filepath.Join(assetsDir, staticCss),
		fonts:  filepath.Join(assetsDir, staticFonts),
		js:     filepath.Join(assetsDir, staticJs),
		img:    filepath.Join(assetsDir, staticImg),
		icons:  filepath.Join(assetsDir, staticIcons),
		pwa:    filepath.Join(assetsDir, staticPwa),
	}
}

func (f *FolderPaths) Assets() string {
	return f.assets
}

func (f *FolderPaths) Static() string {
	return f.static
}

func (f *FolderPaths) Css() string {
	return f.css
}

func (f *FolderPaths) Fonts() string {
	return f.fonts
}

func (f *FolderPaths) Js() string {
	return f.js
}

func (f *FolderPaths) Img() string {
	return f.img
}

func (f *FolderPaths) Icons() string {
	return f.icons
}

func (f *FolderPaths) Pwa() string {
	return f.pwa
}
