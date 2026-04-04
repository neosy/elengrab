package shcode

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"strings"
)

const (
	// Contains all lowercase and uppercase English letters and digits
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// Maximum length for the URL code
	maxShortCodeLength = 40
)

// GenerateShortCode generates a short, URL-safe string based on the SHA-256 hash of the input `key` string.
// This can be used to produce deterministic and compact identifiers (e.g., short codes for links).
//
// Parameters:
//   - key: a string to be hashed (e.g., a UUID, URL, or a combined string like uuid:url:timestamp).
//   - length: an optional maximum length of the generated code. If 0, the full base64-encoded hash (without padding) is returned.
//     The maximum length for the full hash is 43 characters, as base64-encoded SHA-256 (without padding) is 43 characters long.
//     If specified and shorter than the hash length, the result will be truncated to the given length.
//
// Returns:
//   - A URL-safe base64-encoded string (without padding) representing the hash of the input string.
//     If `key` is empty, an empty string is returned.
func GenerateShortCode(key string, length uint8) string {
	// Check that the key is not an empty string
	if key == "" {
		return ""
	}

	// If length is 0, set the default length (maxURLCodeLength)
	if length == 0 {
		length = maxShortCodeLength
	}

	// Hash using SHA-256
	hash := sha256.Sum256([]byte(key))

	// Encode to base64 (URL-safe), remove padding ('=' characters)
	code := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(hash[:])

	// Replace '-' and '_' with random characters
	code = strings.ReplaceAll(code, "-", randomChar())
	code = strings.ReplaceAll(code, "_", randomChar())

	// If a length is specified and it is less than the generated code length
	// Trim the code to the specified length
	if len(code) > int(length) {
		code = code[:length]
	}

	// Trim the code if it exceeds the maximum length
	if len(code) > maxShortCodeLength {
		code = code[:maxShortCodeLength]
	}

	return code
}

// Generate a random character (digit or letter)
func randomChar() string {
	return string(letters[rand.Intn(len(letters))])
}
