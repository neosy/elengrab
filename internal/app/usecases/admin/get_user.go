package admin

import (
	"context"

	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

func (a *Admin) GetAllUsersWithoutGuest(ctx context.Context) ([]*dauth.User, error) {
	return a.adminPanel.GetAllUsersWithoutGuest(ctx)
}
