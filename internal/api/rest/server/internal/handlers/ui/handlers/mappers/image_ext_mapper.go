package mappers

import "github.com/neosy/elengrab/pkg/httpx"

func (m *Mappers) MapImageExtToContentType(ext string) string {
	contentType, exists := httpx.ImageContentTypeFromExt(ext)
	if !exists {
		contentType = "application/octet-stream"
	}
	return contentType
}
