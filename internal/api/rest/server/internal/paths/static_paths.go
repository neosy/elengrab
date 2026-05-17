package httppaths

const (
	// Groups
	GroupStatic = "/static"

	GroupStaticCss             = GroupStatic + "/css"
	GroupStaticFonts           = GroupStatic + "/fonts"
	GroupStaticJs              = GroupStatic + "/js"
	GroupStaticImages          = GroupStatic + "/images"
	GroupStaticIcons           = GroupStatic + "/icons"
	GroupStaticPwa             = GroupStatic + "/pwa"
	GroupStaticThumbnails      = GroupStatic + "/thumbnails"
	GroupStaticYoutubeChannels = GroupStatic + "/ytchannels"

	// Path files
	PathCssFiles        = "/css/{filepath:*}"
	PathFontFiles       = "/fonts/{filepath:*}"
	PathImageFiles      = "/images/{filepath:*}"
	PathImageFaviconICO = "/images/favicon.ico"
	PathIconFiles       = "/icons/{filepath:*}"
	PathJsFiles         = "/js/{filepath:*}"
	PathPwaFiles        = "/pwa/{filepath:*}"
	PathThumbnail       = "/thumbnails/{thumbnailId}"
	PathYoutubeChannel  = "/ytchannels/{channelId}"
)
