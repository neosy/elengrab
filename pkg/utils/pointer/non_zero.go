package uptr

import "github.com/google/uuid"

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
