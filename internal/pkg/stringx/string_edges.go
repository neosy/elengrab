package stringx

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Capitalize returns a copy of the string s with the first letter capitalized.
func Capitalize(s string) string {
	if s == "" {
		return s
	}

	// Get the first rune and its size in bytes
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}

// LowerFirst returns a copy of the string s with the first letter in lowercase.
func LowerFirst(s string) string {
	if s == "" {
		return ""
	}

	// Get the first rune and its size in bytes
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[size:]
}

// LowerFirstWord lowers the first word of the string.
func LowerFirstWord(s string) string {
	if s == "" {
		return ""
	}

	words := strings.SplitN(s, " ", 2)
	words[0] = strings.ToLower(words[0])

	return strings.Join(words, " ")
}
