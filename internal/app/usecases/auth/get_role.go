package auth

import (
	"context"

	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

func (a *Auth) GetAllRoles(ctx context.Context) ([]*dauth.Role, error) {
	return a.role.GetAll(ctx)
}

func (a *Auth) GetAllRolesWithoutGuest(ctx context.Context) ([]*dauth.Role, error) {
	return a.role.GetAllWithoutGuest(ctx)
}
