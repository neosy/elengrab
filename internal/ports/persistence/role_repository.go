package persistence

import (
	"context"

	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

type RoleRepository interface {
	Transactional
	Insert(ctx context.Context, role *dauth.Role) error
	Update(ctx context.Context, role *dauth.Role) error
	Save(ctx context.Context, role *dauth.Role) error

	Find(ctx context.Context, roleID string) (*dauth.Role, error)
	Exists(ctx context.Context, roleID string) (bool, error)

	GetAll(ctx context.Context) ([]*dauth.Role, error)

	WithoutGuest() RoleRepository
}
