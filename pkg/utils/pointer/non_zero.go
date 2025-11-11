package uptr

import (
	"reflect"

	"github.com/google/uuid"
)

// NonZeroString returns a pointer to v, or nil if v is empty.
func NonZeroString(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

// NonZeroUuid returns a pointer to v, or nil if v is zero.
func NonZeroUuid(v uuid.UUID) *uuid.UUID {
	if v == uuid.Nil {
		return nil
	}
	return &v
}

// NonZero returns a pointer to v if it is not the zero value of its type,
// otherwise returns nil.
func NonZero[T any](v T) *T {
	// Check if v is the zero value for its type
	if reflect.ValueOf(v).IsZero() {
		return nil
	}
	return &v
}
