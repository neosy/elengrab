package auth

import (
	"context"

	autherr "github.com/neosy/elengrab/internal/app/usecases/auth/errors"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
)

func (a *Auth) ValidateSession(ctx context.Context, token string) (*dto.AuthUserResponse, error) {
	session, err := a.userSession.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, autherr.ErrSessionNotFound
	}

	if session.Expired() {
		return nil, autherr.ErrSessionExpired
	}

	user, err := a.user.FindByUserID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		a.logger.Debug("User not found", "token", token)
		return nil, autherr.ErrUserNotFound
	}

	if !user.IsActive {
		a.logger.Info("User is not active", "userID", user.UserID)
		return nil, autherr.ErrUserIsNotActive
	}

	if user.DeletedAt != nil {
		a.logger.Info("User deleted", "userID", user.UserID)
		return nil, autherr.ErrUserDeleted
	}

	userCtx := a.mappers.MapUserSessionDomainToUserContext(
		user,
		session,
		a.sessionRefreshPredicate(),
	)

	return userCtx, nil
}
