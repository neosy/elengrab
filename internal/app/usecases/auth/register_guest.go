package auth

import (
	"context"

	authdto "github.com/neosy/elengrab/internal/app/usecases/auth/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

func (a *Auth) RegisterGuest(ctx context.Context) (*authdto.UserContext, error) {
	var (
		user    *dauth.User
		session *dauth.UserSession
	)

	createUser := func(ctx context.Context) error {
		var err error
		user, err = a.createGuest(ctx)
		if err != nil {
			return err
		}

		session, err = a.createSession(ctx, user.UserID)
		if err != nil {
			return err
		}
		return nil
	}

	err := a.user.Tx(ctx, createUser)
	if err != nil {
		return nil, err
	}

	userCtx := a.mappers.MapUserSessionDomainToUserContext(
		user,
		session,
		a.sessionRefreshPredicate(),
	)

	return userCtx, nil
}

func (u *Auth) createGuest(ctx context.Context) (*dauth.User, error) {
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
