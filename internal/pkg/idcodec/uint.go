package idcodec

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"unsafe"
)

// EncodeUintBase64URL converts an unsigned integer into a URL-safe short string representation.
// The encoded size depends on the size of the integer type.
func EncodeUintBase64URL[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](value T) string {
	size := unsafe.Sizeof(value)
	data := make([]byte, size)

	switch size {
	case 1:
		data[0] = byte(value)
	case 2:
		binary.BigEndian.PutUint16(data, uint16(value))
	case 4:
		binary.BigEndian.PutUint32(data, uint32(value))
	case 8:
		binary.BigEndian.PutUint64(data, uint64(value))
	}

	return base64.RawURLEncoding.EncodeToString(data)
}

// DecodeUintBase64URL converts a URL-safe short string representation into an unsigned integer.
func DecodeUintBase64URL[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](value string) (T, error) {
	var result T
	size := unsafe.Sizeof(result)

	data, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return 0, err
	}

	if uintptr(len(data)) != size {
		return 0, fmt.Errorf("invalid uint data length: %d, expected: %d", len(data), size)
	}

	switch size {
	case 1:
		result = T(data[0])
	case 2:
		result = T(binary.BigEndian.Uint16(data))
	case 4:
		result = T(binary.BigEndian.Uint32(data))
	case 8:
		result = T(binary.BigEndian.Uint64(data))
	}

	return result, nil
}
