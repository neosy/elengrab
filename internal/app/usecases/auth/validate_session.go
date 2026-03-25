package auth

import (
	"context"

	authdto "github.com/neosy/elengrab/internal/app/usecases/auth/dto"
)

func (a *Auth) ValidateSession(ctx context.Context, token string) (*authdto.UserContext, error) {
	session, err := a.userSession.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, ErrSessionNotFound
	}

	if session.Expired() {
		return nil, ErrSessionExpired
	}

	user, err := a.user.FindByUserID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	userCtx := a.mappers.MapUserSessionDomainToUserContext(
		user,
		session,
		a.sessionRefreshPredicate(),
	)

	return userCtx, nil
}
