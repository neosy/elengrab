package dauth

import "time"

type Role struct {
	// Unique role identifier (can be a readable key like 'admin', 'guest')
	RoleID string

	// Human-readable role name, must be unique across the system
	Name string

	// Optional description of role
	Description *string

	// Timestamp when the record was created
	CreatedAt time.Time

	// Timestamp when the record was last updated
	UpdatedAt time.Time
}
