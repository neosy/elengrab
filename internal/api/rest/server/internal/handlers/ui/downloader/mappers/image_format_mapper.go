package mappers

var (
	mapImageFormatToContentType = map[string]string{
		"jpg":  "image/jpeg",
		"jpeg": "image/jpeg",
		"png":  "image/png",
		"webp": "image/webp",
		"svg":  "image/svg+xml",
		"ico":  "image/x-icon",
	}
)

func (m *Mappers) MapImageFormatToContentType(imageFormat string) string {
	contentType := mapImageFormatToContentType[imageFormat]
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return contentType
}
