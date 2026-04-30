package httppaths

const (
	// Groups
	GroupStatic = "/static"

	GroupStaticCss            = GroupStatic + "/css"
	GroupStaticJs             = GroupStatic + "/js"
	GroupStaticImg            = GroupStatic + "/img"
	GroupStaticIcon           = GroupStatic + "/icon"
	GroupStaticPwa            = GroupStatic + "/pwa"
	GroupStaticThumbnail      = GroupStatic + "/thumbnail"
	GroupStaticYoutubeChannel = GroupStatic + "/ytchannel"

	// Path files
	PathCssFiles       = "/css/{filepath:*}"
	PathImgFiles       = "/img/{filepath:*}"
	PathImgFaviconICO  = "/img/favicon.ico"
	PathIconFiles      = "/icon/{filepath:*}"
	PathJsFiles        = "/js/{filepath:*}"
	PathPwaFiles       = "/pwa/{filepath:*}"
	PathThumbnail      = "/thumbnail/{thumbnailId}"
	PathYoutubeChannel = "/ytchannel/{channelId}"
)
