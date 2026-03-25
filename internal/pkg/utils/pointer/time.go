// Package uptr provides utility functions for working with pointers.
// It includes functions to return pointers to different time types.
package uptr

import "time"

// TimeDuration returns a pointer to the provided time Duration.
func TimeDuration(t time.Duration) *time.Duration {
	return &t
}
