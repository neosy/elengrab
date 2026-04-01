package dauth

import (
	"time"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type User struct {
	// Unique user identifier (UUID)
	UserID uuid.UUID

	// Username/login, must be unique
	Login dtypes.Login

	// User email address, optional
	Email *string

	// User password hash, optional
	PasswordHash *string

	// Timestamp when the user's password was last updated
	PasswordUpdatedAt *time.Time

	// Active status
	IsActive bool

	// Timestamp when the record was created
	CreatedAt time.Time

	// Timestamp when the record was last updated
	UpdatedAt time.Time

	// Timestamp when the record was soft deleted
	DeletedAt *time.Time

	// Roles
	Roles []dtypes.UserRole
}
