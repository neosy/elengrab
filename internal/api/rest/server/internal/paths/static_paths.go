package httppaths

const (
	// Groups
	GroupStatic = "/static"

	// Path files
	PathCssFiles = "/css/{filepath:*}"
	PathImgFiles = "/img/{filepath:*}"
	PathJsFiles  = "/js/{filepath:*}"
	PathPwaFiles = "/pwa/{filepath:*}"
)
