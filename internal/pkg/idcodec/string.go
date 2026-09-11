package idcodec

import (
	"encoding/base64"
)

// EncodeStringBase64URL converts string into a URL-safe short string representation.
func EncodeStringBase64URL(v string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(v))
}

// DecodeStringBase64URL converts a URL-safe encoded string back into the original string.
func DecodeStringBase64URL(v string) (string, error) {
	data, err := base64.RawURLEncoding.DecodeString(v)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
