package admin

import (
	"context"

	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

func (a *admin) GetAllUsersWithoutGuest(ctx context.Context) ([]*dauth.User, error) {
	return a.adminPanel.GetAllUsersWithoutGuest(ctx)
}
