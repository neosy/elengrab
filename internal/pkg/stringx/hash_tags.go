package stringx

import (
	"regexp"
	"strings"
)

var (
	hashtagRegex         = regexp.MustCompile(`#\S+`)
	trailingHashtagRegex = regexp.MustCompile(`(?:\s*#[\p{L}\p{N}_]+)+$`)
)

// RemoveHashtags removes hashtags from text.
func RemoveHashtags(s string) string {
	s = hashtagRegex.ReplaceAllString(s, "")

	return strings.Join(strings.Fields(s), " ")
}

// RemoveTrailingHashtags removes hashtags only from the end of text.
func RemoveTrailingHashtags(s string) string {
	return strings.TrimSpace(trailingHashtagRegex.ReplaceAllString(s, ""))
}
