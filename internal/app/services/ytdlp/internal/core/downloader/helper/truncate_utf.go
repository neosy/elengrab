package helper

import "unicode/utf8"

// TruncateUTF8 truncates the given UTF-8 string to a maximum of maxBytes bytes.
// If the string is already within the byte limit, it returns the original string.
// Otherwise, it finds the last valid UTF-8 rune before exceeding the byte limit and truncates there.
func TruncateUTF8(s string, maxBytes int) string {
	if len(s) <= maxBytes {
		return s
	}
	// Truncates the string `s` to ensure it is a valid UTF-8 encoded string without exceeding `maxBytes`.
	for maxBytes > 0 && !utf8.ValidString(s[:maxBytes]) {
		maxBytes--
	}
	return s[:maxBytes]
}
