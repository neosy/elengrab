package dlink

import (
	"time"

	"github.com/google/uuid"
)

type LinkClick struct {
	// Unique identifier for the click event
	LinkClickID uuid.UUID

	// ID of the link that was clicked
	LinkID uuid.UUID

	// The IP address from which the link was accessed
	IPAddress string

	// Full short URL, including domain, e.g., "https://s.nhub.ru/abc123"
	ShortURL string

	// User ID who clicked the link (nullable, if not logged in or unknown)
	ClickedBy *string

	// Timestamp of the click event
	ClickedAt time.Time

	// User agent or device info (optional for tracking purposes)
	UserAgent *string

	// Referrer URL (optional, if available)
	Referrer *string

	// Timestamp when the event was created
	CreatedAt time.Time
}
