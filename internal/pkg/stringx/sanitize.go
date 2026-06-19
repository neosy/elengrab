package stringx

import (
	"regexp"
	"strings"
)

var (
	sanitizeForMetaURLRegex    = regexp.MustCompile(`https?://\S+`)
	sanitizeForMetaTrimSymbols = " :-|•"
)

// SanitizeForMetaPreview cleans text for Open Graph usage.
// - removes URLs
// - normalizes whitespace
// - trims result
func SanitizeForMetaPreview(s string, max int, suffix string) string {
	lines := strings.Split(s, "\n")

	var cleaned []string

	for i, line := range lines {
		// trim spaces
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		// remove URLs
		if i == 0 {
			line = sanitizeForMetaURLRegex.ReplaceAllString(line, "")
		}

		if i > 0 && sanitizeForMetaURLRegex.MatchString(line) {
			continue
		}

		// trim spaces
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		cleaned = append(cleaned, line)
	}

	result := strings.Join(cleaned, " ")
	result = strings.Trim(result, sanitizeForMetaTrimSymbols)

	return TruncateWithSuffix(result, max, suffix)
}
