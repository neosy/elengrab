// Package uptr provides utility functions for working with pointers.
// It includes functions to return pointers to different any types.
package uptr

// Any returns a pointer to the provided type.
func Any[T any](v T) *T {
	return &v
}
