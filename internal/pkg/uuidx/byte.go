package uuidx

import (
	"fmt"
	"uuid"
)

// FromBytes converts a byte slice to a UUID.
func FromBytes(b []byte) (uuid.UUID, error) {
	var id uuid.UUID

	if len(b) != 16 {
		return id, fmt.Errorf("invalid UUID (got %d bytes)", len(b))
	}

	copy(id[:], b)

	return id, nil
}
