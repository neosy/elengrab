package idcodec

import (
	"encoding/base64"

	"github.com/google/uuid"
)

// EncodeUUIDBase64URL converts UUID into a URL-safe short string representation.
// The encoded value can be decoded back to the original UUID using DecodeUUID.
func EncodeUUIDBase64URL(id uuid.UUID) string {
	return base64.RawURLEncoding.EncodeToString(id[:])
}

// DecodeUUIDBase64URL converts a URL-safe encoded string back into the original UUID.
// Returns an error if the value is not a valid encoded UUID.
func DecodeUUIDBase64URL(value string) (uuid.UUID, error) {
	data, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.FromBytes(data)
}
