// Package uptr provides utility functions for working with pointers.
// It includes functions to return pointers to different string types.
package uptr

// String returns a pointer to the provided string.
func String(s string) *string {
	return &s
}
