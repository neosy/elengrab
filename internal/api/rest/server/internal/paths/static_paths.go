package httppaths

import (
	"strings"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
)

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

func BuildThumbnailPath(thumbID uuid.UUID) string {
	id := idcodec.EncodeUUIDBase64URL(thumbID)
	return GroupStatic + strings.Replace(PathThumbnail, "{thumbnailId}", id, 1)
}

func BuildImagePath(fileName string) string {
	return GroupStatic + strings.Replace(PathImageFiles, "{filepath:*}", fileName, 1)
}

func BuildIconPath(fileName string) string {
	return GroupStatic + strings.Replace(PathIconFiles, "{filepath:*}", fileName, 1)
}
