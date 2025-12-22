package nfasthttp

import (
	"regexp"
	"strings"
)

// SanitizeFileName normalizes a file name to be safe for filesystem storage and HTTP sendfile usage
func SanitizeFileName(name string) string {
	// Trim spaces at the beginning and end
	safe := strings.TrimSpace(name)

	// Replace all problematic characters with safe alternatives
	replacer := strings.NewReplacer(
		"#", "--",
		"?", "_",
		"%", "_",
		"\\", "-",
		"/", "-",
		":", "-",
		"|", "-",
		"\"", "",
		"<", "",
		">", "",
	)
	safe = replacer.Replace(safe)

	// Replace control characters with underscore
	safe = regexp.MustCompile(`[\x00-\x1F]`).ReplaceAllString(safe, "_")

	// Normalize spaces and underscores
	safe = regexp.MustCompile(`\s+`).ReplaceAllString(safe, " ")
	safe = regexp.MustCompile(`_+`).ReplaceAllString(safe, "_")

	// Remove space before extension
	safe = regexp.MustCompile(`\s+(\.)`).ReplaceAllString(safe, "$1")

	// UTF-8 safe length limit (255 bytes typical FS limit)
	safe = truncateUTF8(safe, 255)

	return safe
}

func truncateUTF8(s string, maxBytes int) string {
	if len(s) <= maxBytes {
		return s
	}

	b := []byte(s)
	if len(b) <= maxBytes {
		return s
	}

	// cut safely at rune boundary
	for maxBytes > 0 && (b[maxBytes]&0xC0) == 0x80 {
		maxBytes--
	}
	return string(b[:maxBytes])
}
