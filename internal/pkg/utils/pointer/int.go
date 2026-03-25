// Package uptr provides utility functions for working with pointers.
// It includes functions to return pointers to different integer types.
package uptr

// Int returns a pointer to the provided int.
func Int(v int) *int {
	return &v
}

// Int8 returns a pointer to the provided int8.
func Int8(v int8) *int8 {
	return &v
}

// Int8 returns a pointer to the provided int16.
func Int16(v int16) *int16 {
	return &v
}

// Int32 returns a pointer to the provided int32.
func Int32(v int32) *int32 {
	return &v
}

// Int64 returns a pointer to the provided int64.
func Int64(v int64) *int64 {
	return &v
}
