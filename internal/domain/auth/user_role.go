package dauth

import (
	"time"

	"github.com/google/uuid"
)

type UserRole struct {
	// Reference to the user (many-to-many relationship)
	UserID uuid.UUID

	// Reference to the role assigned to the user
	RoleID string

	// Timestamp when the record was created
	CreatedAt time.Time
}
