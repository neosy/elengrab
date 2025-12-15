package mappers

var (
	mapImageFormatToContentType = map[string]string{
		"jpg":  "image/jpeg",
		"jpeg": "image/jpeg",
		"png":  "image/png",
		"webp": "image/webp",
	}
)

func (m *Mappers) MapImageFormatToContentType(imageFormat string) string {
	contentType := mapImageFormatToContentType[imageFormat]
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return contentType
}
