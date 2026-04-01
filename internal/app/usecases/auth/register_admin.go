package auth

import (
	"context"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (a *Auth) RegisterAdmin(ctx context.Context, req *dto.RegisterAdminRequest) (*dto.AuthUserResponse, error) {
	if req == nil {
		return nil, errorx.New("function parameter is nil", exceptionx.ERROR)
	}

	var (
		user         *dauth.User
		session      *dauth.UserSession
		passwordHash *string
	)

	if req.Password != "" {
		hash, err := a.hashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		passwordHash = &hash
	}

	createUser := func(ctx context.Context) error {
		var err error
		user, err = a.createAdmin(ctx, req.Login, req.Email, passwordHash)
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
		nil,
	)

	return userCtx, nil
}

func (u *Auth) createAdmin(
	ctx context.Context,
	login string,
	email string,
	passwordHash *string,
) (*dauth.User, error) {
	userID, err := u.user.CreateAdmin(ctx, dtypes.NewLogin(login), email, passwordHash)
	if err != nil {
		return nil, err
	}

	user, err := u.user.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
