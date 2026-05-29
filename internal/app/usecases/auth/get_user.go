package auth

import (
	"context"

	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (a *Auth) FindByUserID(ctx context.Context, userID uuid.UUID) (*dauth.User, error) {
	return a.user.FindByUserID(ctx, userID)
}

func (a *Auth) ExistsUserByLogin(ctx context.Context, login string) (bool, error) {
	return a.user.ExistsByLogin(ctx, dtypes.NewLogin(login))
}

func (a *Auth) GetAllUsers(ctx context.Context) ([]*dauth.User, error) {
	return a.user.GetAllUsers(ctx)
}

func (a *Auth) GetAllUsersWithoutGuest(ctx context.Context) ([]*dauth.User, error) {
	return a.user.GetAllUsersWithoutGuest(ctx)
}
