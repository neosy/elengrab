package filex

import (
	"regexp"
	"strings"
)

// SanitizeFileName replaces characters not allowed in file names with "_"
func SanitizeFileName(name string) string {
	// Trim spaces at the beginning and end
	safe := strings.TrimSpace(name)

	// Replace all invalid characters with underscore
	// Windows forbidden: \ / : * ? " < > |
	// Unix forbidden: /
	safe = regexp.MustCompile(`[:|]`).ReplaceAllString(safe, "-")
	safe = regexp.MustCompile(`[/\\]`).ReplaceAllString(safe, "-")
	safe = regexp.MustCompile(`["<>]`).ReplaceAllString(safe, "")
	safe = regexp.MustCompile(`[?*\x00-\x1F]`).ReplaceAllString(safe, "_")

	// Optional: replace multiple underscores with a single one
	safe = regexp.MustCompile(`_+`).ReplaceAllString(safe, "_")

	// Optional: limit length (e.g., 255 characters)
	if len(safe) > 255 {
		safe = safe[:255]
	}

	return safe
}
