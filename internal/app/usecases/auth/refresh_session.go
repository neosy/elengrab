package auth

import (
	"context"
	"time"

	authdto "github.com/neosy/elengrab/internal/app/usecases/auth/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

func (u *Auth) shouldRefreshSession(expiresAt time.Time) bool {
	return time.Until(expiresAt) <= sessionRefreshInterval
}

func (a *Auth) sessionRefreshPredicate() func(*dauth.UserSession) bool {
	return func(session *dauth.UserSession) bool {
		return a.shouldRefreshSession(session.ExpiresAt)
	}
}

func (a *Auth) RefreshSession(ctx context.Context, token string) (*authdto.UserContext, error) {
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
			return ErrSessionNotFound
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
