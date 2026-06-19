package stringx

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	truncateTrimSymbols = " :-|•"
)

// Truncate cuts string by rune length.
// It is safe for UTF-8 (works with Cyrillic, emojis).
func Truncate(s string, max int) string {
	if max <= 0 {
		return ""
	}

	r := []rune(s)

	if len(r) <= max {
		return s
	}

	return string(r[:max])
}

// Truncate cuts string by rune length and appends suffix if needed.
// It is safe for UTF-8 (works with Cyrillic, emojis).
func TruncateWithSuffix(s string, max int, suffix string) string {
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

// TruncateWords cuts string by rune length without breaking words.
// It removes the last word only if it was partially truncated.
// It is safe for UTF-8.
func TruncateWords(s string, max int) string {
	if max <= 0 {
		return ""
	}

	r := []rune(s)

	if len(r) <= max {
		return s
	}

	truncated := string(r[:max])

	// If truncation happened in the middle of a word, remove that word.
	if !unicode.IsSpace(r[max-1]) && !unicode.IsSpace(r[max]) {
		if index := strings.LastIndexAny(truncated, " \t\n"); index > 0 {
			return strings.TrimSpace(truncated[:index])
		}
	}

	return strings.TrimSpace(truncated)
}

// TruncateBytes truncates the given UTF-8 string to a maximum of maxBytes bytes.
// If the string is already within the byte limit, it returns the original string.
// Otherwise, it finds the last valid UTF-8 rune before exceeding the byte limit and truncates there.
func TruncateBytes(s string, maxBytes int) string {
	if len(s) <= maxBytes {
		return s
	}
	// Truncates the string `s` to ensure it is a valid UTF-8 encoded string without exceeding `maxBytes`.
	for maxBytes > 0 && !utf8.ValidString(s[:maxBytes]) {
		maxBytes--
	}
	return s[:maxBytes]
}

// TruncateBytesWords truncates a UTF-8 string to a maximum of maxBytes bytes.
// It avoids breaking words when possible.
// If the first word exceeds maxBytes, it truncates at a valid UTF-8 boundary.
func TruncateBytesWords(s string, maxBytes int) string {
	if maxBytes <= 0 {
		return ""
	}

	if len(s) <= maxBytes {
		return s
	}

	runes := []rune(s)

	bytesCount := 0
	lastSpace := -1
	cut := 0

	for i, r := range runes {
		runeBytes := len(string(r))

		if bytesCount+runeBytes > maxBytes {
			break
		}

		bytesCount += runeBytes
		cut = i + 1

		if unicode.IsSpace(r) {
			lastSpace = i
		}
	}

	if cut == len(runes) {
		return s
	}

	// If truncation happened in the middle of a word, cut to the previous word.
	if cut < len(runes) && !unicode.IsSpace(runes[cut-1]) && !unicode.IsSpace(runes[cut]) {
		if lastSpace >= 0 {
			cut = lastSpace
		}
	}

	return strings.TrimSpace(string(runes[:cut]))
}
