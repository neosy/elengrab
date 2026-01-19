package auth

import (
	"context"

	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

func (u *Auth) CreateGuestUser(ctx context.Context) (*dauth.User, error) {
	userID, err := u.user.CreateGuest(ctx)
	if err != nil {
		return nil, err
	}

	user, err := u.user.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
