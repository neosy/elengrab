package auth

import (
	"context"
	"time"

	autherr "github.com/neosy/elengrab/internal/app/usecases/auth/errors"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

func (a *auth) RefreshSession(ctx context.Context, token string) (*dto.AuthUserResponse, error) {
	var (
		user       *dauth.User
		newSession *dauth.UserSession
	)

	refresh := func(ctx context.Context) error {
		session, err := a.userSession.FindByToken(ctx, token)
		if err != nil {
			return err
		}
		if session == nil {
			return autherr.ErrSessionNotFound
		}

		if !a.shouldRefreshSession(session.ExpiresAt) {
			newSession = session
		}

		if newSession == nil {
			newSession, err = a.createSession(ctx, session.UserID)
			if err != nil {
				return err
			}
		}

		user, err = a.user.GetByUserID(ctx, newSession.UserID)
		if err != nil {
			return err
		}

		return nil
	}

	err := a.userSession.Tx(ctx, refresh)
	if err != nil {
		return nil, err
	}

	userCtx := a.mappers.MapUserSessionDomainToUserContext(
		user,
		newSession,
		a.sessionRefreshPredicate(),
	)

	return userCtx, nil
}

func (u *auth) shouldRefreshSession(expiresAt time.Time) bool {
	return time.Until(expiresAt) <= u.sessionRefreshInterval
}

func (a *auth) sessionRefreshPredicate() func(*dauth.UserSession) bool {
	return func(session *dauth.UserSession) bool {
		return a.shouldRefreshSession(session.ExpiresAt)
	}
}
