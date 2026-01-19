package dauth

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	// Unique user identifier (UUID)
	UserID uuid.UUID

	// Username/login, must be unique
	Login string

	// User email address, optional
	Email *string

	// Indicates guest status
	IsGuest bool

	// Active status
	IsActive bool

	// Timestamp when the record was created
	CreatedAt time.Time

	// Timestamp when the record was last updated
	UpdatedAt time.Time
}
