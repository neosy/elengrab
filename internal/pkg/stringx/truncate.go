package stringx

import "strings"

var (
	truncateTrimSymbols = " :-|•"
)

// Truncate cuts string by rune length and appends suffix if needed.
// It is safe for UTF-8 (works with Cyrillic, emojis).
func Truncate(s string, max int, suffix string) string {
	if max <= 0 {
		return ""
	}

	r := []rune(s)

	if len(r) <= max {
		return s
	}

	if max < len([]rune(suffix)) {
		return string(r[:max])
	}

	line := string(r[:max-len([]rune(suffix))])
	line = strings.Trim(line, truncateTrimSymbols)

	return line + suffix
}
