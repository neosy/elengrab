package shcode

import (
	"crypto/sha256"
	"encoding/base64"
	"hash/fnv"
	"math/rand"
	"strings"
)

const (
	// Contains all lowercase and uppercase English letters and digits
	base62 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// Maximum length for the URL code
	maxShortCodeLength = 40
	// Default length for the URL code
	defaultShortCodeLength = 8
)

// GenerateShortCode generates a short, URL-safe string based on the SHA-256 hash of the input `key` string.
// This can be used to produce deterministic and compact identifiers (e.g., short codes for links).
//
// Parameters:
//   - key: a string to be hashed (e.g., a UUID, URL, or a combined string like uuid:url:timestamp).
//   - length: an optional maximum length of the generated code. If 0, the full base64-encoded hash (without padding) is returned.
//     The maximum length for the full hash is 43 characters, as base64-encoded SHA-256 (without padding) is 43 characters long.
//     If specified and shorter than the hash length, the result will be truncated to the given length.
//   - deterministic - determines whether the generated code is deterministic;
//     when true, the same key always produces the same code
//
// Returns:
//   - A URL-safe base64-encoded string (without padding) representing the hash of the input string.
//     If `key` is empty, an empty string is returned.
func GenerateShortCode(key string, length uint8, deterministic bool) string {
	// Check that the key is not an empty string
	if key == "" {
		return ""
	}

	// If length is 0, set the default length (defaultShortCodeLength)
	if length == 0 {
		length = defaultShortCodeLength
	}

	var code string
	if length <= 5 {
		code = string(hashFnv32ToBase62(key, length))
	} else if length <= 8 {
		code = string(hashFnv64ToBase62(key, length))
	} else {
		code = hashSha256ToBase62(key, length, deterministic)
	}

	return code
}

// Generate a random character (digit or letter)
func randomChar() string {
	return string(base62[rand.Intn(len(base62))])
}

// hashFnv32ToBase62
// Length up to 6 characters
func hashFnv32ToBase62(key string, length uint8) string {
	const maxLength = 6

	// FNV-1a 32-bit hash
	h := fnv.New32a() // hash.Hash32
	h.Write([]byte(key))
	num := h.Sum32()

	code := ""
	if num == 0 {
		code = string(base62[0])
	} else {
		for num > 0 {
			code = string(base62[num%62]) + code
			num /= 62
		}
	}

	l := int(length)

	if l > maxLength {
		l = maxLength
	}

	// Pad with '0' if shorter than length
	for len(code) < l {
		code = string(base62[0]) + code
	}

	// Trim if longer than length
	if len(code) > l {
		code = code[:l]
	}

	return code
}

// hashFnv64ToBase62
// Length up to 11 characters
func hashFnv64ToBase62(key string, length uint8) string {
	const maxLength = 11

	// FNV-1a 64-bit hash
	h := fnv.New64a() // hash.Hash64
	h.Write([]byte(key))
	num := h.Sum64()

	code := ""
	if num == 0 {
		code = string(base62[0])
	} else {
		for num > 0 {
			code = string(base62[num%62]) + code
			num /= 62
		}
	}

	l := int(length)

	if l > maxLength {
		l = maxLength
	}

	// Pad with '0' if shorter than length
	for len(code) < l {
		code = string(base62[0]) + code
	}

	// Trim if longer than length
	if len(code) > l {
		code = code[:l]
	}

	return code
}

// hashSha256ToBase62
// Length up to 40 characters
func hashSha256ToBase62(key string, length uint8, deterministic bool) string {
	// Hash using SHA-256
	hash := sha256.Sum256([]byte(key))

	// Encode to base64 (URL-safe), remove padding ('=' characters)
	code := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(hash[:])

	// Replace '-' and '_' with random characters
	if deterministic {
		code = strings.ReplaceAll(code, "-", "")
		code = strings.ReplaceAll(code, "_", "")
	} else {
		// Можно добавить случайные символы для не детерминированного варианта
		code = strings.ReplaceAll(code, "-", randomChar())
		code = strings.ReplaceAll(code, "_", randomChar())
	}

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
