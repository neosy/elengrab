package httppaths

const (
	// Groups
	GroupStatic = "/static"

	// Path files
	PathCssFiles      = "/css/{filepath:*}"
	PathImgFiles      = "/img/{filepath:*}"
	PathImgFaviconICO = "/img/favicon.ico"
	PathIconFiles     = "/icon/{filepath:*}"
	PathJsFiles       = "/js/{filepath:*}"
	PathPwaFiles      = "/pwa/{filepath:*}"
)
