package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	autherr "github.com/neosy/elengrab/internal/app/usecases/auth/errors"
	authtoken "github.com/neosy/elengrab/internal/app/usecases/auth/internal/token"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (u *Auth) createSession(ctx context.Context, userID uuid.UUID) (*dauth.UserSession, error) {
	token, err := authtoken.GenerateToken(authtoken.CookieToken)
	if err != nil {
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	session := &dauth.UserSession{
		UserID:       userID,
		SessionToken: token,
		ExpiresAt:    time.Now().UTC().Add(sessionTTL),
	}

	sessionID, err := u.userSession.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	session, err = u.userSession.FindBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, autherr.ErrSessionNotFound
	}

	return session, nil
}
