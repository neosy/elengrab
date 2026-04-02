package auth

import (
	"context"

	autherr "github.com/neosy/elengrab/internal/app/usecases/auth/errors"
	idto "github.com/neosy/elengrab/internal/app/usecases/auth/internal/dto"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (a *Auth) RegisterAdmin(ctx context.Context, req *dto.RegisterAdminRequest) (*dto.AuthUserResponse, error) {
	if req == nil {
		return nil, autherr.ErrFunctionNilParameter
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
		createReq := &idto.CreateUserRequest{
			Login:        dtypes.NewLogin(req.Login),
			Email:        req.Email,
			PasswordHash: passwordHash,
		}
		user, err = a.createAdmin(ctx, createReq)
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

func (u *Auth) createAdmin(ctx context.Context, req *idto.CreateUserRequest) (*dauth.User, error) {
	userID, err := u.user.CreateAdmin(ctx, req)
	if err != nil {
		return nil, err
	}

	user, err := u.user.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
