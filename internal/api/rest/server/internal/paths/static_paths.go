package httppaths

const (
	// Groups
	GroupStatic = "/static"

	GroupCss  = GroupStatic + "/css"
	GroupJs   = GroupStatic + "/js"
	GroupImg  = GroupStatic + "/img"
	GroupIcon = GroupStatic + "/icon"
	GroupPwa  = GroupStatic + "/pwa"

	// Path files
	PathCssFiles      = "/css/{filepath:*}"
	PathImgFiles      = "/img/{filepath:*}"
	PathImgFaviconICO = "/img/favicon.ico"
	PathIconFiles     = "/icon/{filepath:*}"
	PathJsFiles       = "/js/{filepath:*}"
	PathPwaFiles      = "/pwa/{filepath:*}"
)
