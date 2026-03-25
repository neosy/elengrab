package httpx

import (
	"mime"
	"strings"
)

var (
	imageMimeToExt = map[string]string{
		"image/jpeg":               "jpg",
		"image/jpg":                "jpg",
		"image/png":                "png",
		"image/webp":               "webp",
		"image/avif":               "avif",
		"image/tiff":               "tiff",
		"image/gif":                "gif",
		"image/svg+xml":            "svg",
		"image/bmp":                "bmp",
		"image/x-icon":             "ico",
		"image/vnd.microsoft.icon": "ico",
		"image/icon":               "ico",
		"image/heif":               "heif",
		"image/heic":               "heic",
	}
	imageExtToMime map[string]string
)

func init() {
	imageExtToMime = make(map[string]string, len(imageMimeToExt))
	for mimeType, ext := range imageMimeToExt {
		ext = strings.ToLower(ext)
		// if there are several MIME per extension, it will take the first one encountered
		if _, ok := imageExtToMime[ext]; !ok {
			imageExtToMime[ext] = mimeType
		}
	}
}

// ImageExtFromContentType returns the file extension corresponding to a given
// MIME Content-Type. For example, "image/jpeg" -> "jpg".
// It returns the extension and a boolean indicating whether the MIME type is supported.
func ImageExtFromContentType(ct string) (string, bool) {
	mediaType, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return "", false
	}
	ext, ok := imageMimeToExt[mediaType]
	return ext, ok
}

// ImageContentTypeFromExt returns the MIME Content-Type corresponding to a given
// image file extension. For example, "png" -> "image/png".
// The lookup is case-insensitive.
// It returns the MIME type and a boolean indicating whether the extension is supported.
func ImageContentTypeFromExt(ext string) (string, bool) {
	ext = strings.ToLower(ext)
	ct, ok := imageExtToMime[ext]
	return ct, ok
}
