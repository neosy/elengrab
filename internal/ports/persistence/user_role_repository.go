package persistence

import (
	"context"

	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

type UserRoleRepository interface {
	Transactional
	Insert(ctx context.Context, userRole *dauth.UserRole) error
	Find(ctx context.Context, userID uuid.UUID, roleID string) (*dauth.UserRole, error)
	Exists(ctx context.Context, userID uuid.UUID, roleID string) (bool, error)
}
