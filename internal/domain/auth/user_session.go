package dauth

import (
	"time"

	"github.com/google/uuid"
)

type UserSession struct {
	// Unique session identifier (UUID)
	SessionID uuid.UUID

	// Associated user identifier (UUID)
	UserID uuid.UUID

	// Random session token stored in cookie
	SessionToken string

	// Timestamp when the record was created
	CreatedAt time.Time

	// Session expiration timestamp
	ExpiresAt time.Time
}

func (s *UserSession) Expired() bool {
	return time.Now().After(s.ExpiresAt)
}
