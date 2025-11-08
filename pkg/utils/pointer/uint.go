// Package uptr provides utility functions for working with pointers.
// It includes functions to return pointers to different unsigned integer types.
package uptr

// Uint returns a pointer to the provided uint.
func Uint(v uint) *uint {
	return &v
}

// Uint8 returns a pointer to the provided uint8.
func Uint8(v uint8) *uint8 {
	return &v
}

// Uint8 returns a pointer to the provided uint16.
func Uint16(v uint16) *uint16 {
	return &v
}

// Uint32 returns a pointer to the provided uint32.
func Uint32(v uint32) *uint32 {
	return &v
}

// Uint64 returns a pointer to the provided uint64.
func Uint64(v uint64) *uint64 {
	return &v
}
