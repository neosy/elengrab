package httpx

import (
	"mime"
	"path/filepath"
	"strings"
)

// ContentTypeByExt returns the MIME content type for the given file extension.
// The extension may be provided with or without a leading dot.
// Unknown extensions fall back to "application/octet-stream".
func ContentTypeByExt(ext string) string {
	if ext == "" {
		return "application/octet-stream"
	}

	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	if contentType := mime.TypeByExtension(ext); contentType != "" {
		return contentType
	}

	return "application/octet-stream" // fallback for unknown file types
}

// ContentTypeFromPath returns the MIME content type based on the file's path.
func ContentTypeFromPath(path string) string {
	return ContentTypeByExt(filepath.Ext(path))
}

// ExtensionByContentType returns the most common file extension for a given Content-Type.
// It correctly strips any parameters (e.g. "; charset=utf-8") and returns the extension
// without the leading dot.
//
// Examples:
//
//	"image/jpeg"                    → "jpg"
//	"image/jpeg; charset=utf-8"     → "jpg"
//	"text/html; charset=utf-8"      → "html"
//	"image/vnd.microsoft.icon"      → "ico"
//	"application/octet-stream"      → ""
func ExtensionByContentType(contentType string) string {
	if contentType == "" {
		return ""
	}

	// ParseMediaType safely removes parameters like charset, boundary, etc.
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		// Fallback: manually cut off parameters if parsing fails
		if idx := strings.Index(contentType, ";"); idx != -1 {
			mediaType = strings.TrimSpace(contentType[:idx])
		} else {
			mediaType = strings.TrimSpace(contentType)
		}
	}

	exts, err := mime.ExtensionsByType(mediaType)
	if err != nil || len(exts) == 0 {
		return "" // unknown content type
	}

	// Return the first (most conventional) extension without leading dot
	ext := exts[0]
	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}
	return ext
}
