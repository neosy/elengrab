package auth

import (
	"context"

	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

func (u *Auth) FindSessionByToken(ctx context.Context, token string) (*dauth.UserSession, error) {
	session, err := u.userSession.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (u *Auth) GetSessionByToken(ctx context.Context, token string) (*dauth.UserSession, error) {
	session, err := u.userSession.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return session, nil
}
