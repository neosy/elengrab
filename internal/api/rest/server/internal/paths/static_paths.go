package httppaths

import (
	"strings"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
)

const (
	// Groups
	StaticGroup = "/static"

	StaticCssGroup             = StaticGroup + "/css"
	StaticFontsGroup           = StaticGroup + "/fonts"
	StaticJsGroup              = StaticGroup + "/js"
	StaticImagesGroup          = StaticGroup + "/images"
	StaticIconsGroup           = StaticGroup + "/icons"
	StaticPwaGroup             = StaticGroup + "/pwa"
	StaticThumbnailsGroup      = StaticGroup + "/thumbnails"
	StaticYoutubeChannelsGroup = StaticGroup + "/ytchannels"

	// Path files
	CssFilesPath        = "/css/{filepath:*}"
	FontFilesPath       = "/fonts/{filepath:*}"
	ImageFilesPath      = "/images/{filepath:*}"
	ImageFaviconICOPath = "/images/favicon.ico"
	IconFilesPath       = "/icons/{filepath:*}"
	JsFilesPath         = "/js/{filepath:*}"
	PwaFilesPath        = "/pwa/{filepath:*}"
	ThumbnailPath       = "/thumbnails/{thumbnailId}"
	YoutubeChannelPath  = "/ytchannels/{channelId}"
)

func BuildThumbnailPath(thumbID uuid.UUID) string {
	id := idcodec.EncodeUUIDBase64URL(thumbID)
	return StaticGroup + strings.Replace(ThumbnailPath, "{thumbnailId}", id, 1)
}

func BuildImagePath(fileName string) string {
	return StaticGroup + strings.Replace(ImageFilesPath, "{filepath:*}", fileName, 1)
}

func BuildIconPath(fileName string) string {
	return StaticGroup + strings.Replace(IconFilesPath, "{filepath:*}", fileName, 1)
}
