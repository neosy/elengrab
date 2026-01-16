package auth

import (
	"context"
	"time"

	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/pkg/errorx"
	"github.com/neosy/elengrab/pkg/errorx/exceptionx"
)

func (u *Auth) GetSessionByToken(ctx context.Context, token string) (*dauth.UserSession, error) {
	session, err := u.userSession.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, errorx.New("user session not found", exceptionx.UNAUTHORIZED)
	}

	if time.Now().After(session.ExpiresAt) {
		return session, errorx.New("user session expired", exceptionx.UNAUTHORIZED)
	}

	return session, nil
}
