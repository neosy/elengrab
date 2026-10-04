package idcodec

import (
	"encoding/base64"
	"encoding/binary"

	"github.com/cespare/xxhash/v2"
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

// EncodeUUIDShortBase64URL converts UUID into an 11-character URL-safe short string.
// The encoded value is derived from the UUID using xxHash64.
func EncodeUUIDShortBase64URL(id uuid.UUID) string {
	hash := xxhash.Sum64(id[:])

	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], hash)

	return base64.RawURLEncoding.EncodeToString(buf[:])
}

// EncodeUUIDBase64URLBits converts UUID into a URL-safe Base64 string
// containing the specified number of bits, limited to 64 bits.
// If bits exceeds 64, it is limited to 64.
func EncodeUUIDBase64URLBits(id uuid.UUID, bits int) string {
	hash := xxhash.Sum64(id[:])

	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], hash)

	if bits > 64 {
		bits = 64
	}

	bytes := (bits + 7) / 8
	if remainder := bits % 8; remainder != 0 {
		buf[bytes-1] &= byte(0xff << (8 - remainder))
	}

	return base64.RawURLEncoding.EncodeToString(buf[:bytes])
}
