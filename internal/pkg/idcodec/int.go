package idcodec

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
)

// EncodeInt64Base64URL converts int64 into a URL-safe short string representation.
func EncodeInt64Base64URL(value int64) string {
	var data [8]byte
	binary.BigEndian.PutUint64(data[:], uint64(value))
	return base64.RawURLEncoding.EncodeToString(data[:])
}

// DecodeInt64Base64URL converts a URL-safe short string representation into int64.
func DecodeInt64Base64URL(value string) (int64, error) {
	data, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return 0, err
	}

	if len(data) != 8 {
		return 0, fmt.Errorf("invalid int64 data length: %d", len(data))
	}

	return int64(binary.BigEndian.Uint64(data)), nil
}
