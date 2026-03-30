package auth

import (
	"context"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (a *Auth) RegisterUser(
	ctx context.Context,
	req *dto.RegisterUserRequest,
) (*dto.AuthUserResponse, error) {
	if req == nil {
		return nil, errorx.New("function parameter is nil", exceptionx.ERROR)
	}

	var (
		user         *dauth.User
		session      *dauth.UserSession
		passwordHash *string
		roles        []dtypes.UserRole
	)

	if req.Password != "" {
		hash, err := a.hashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		passwordHash = &hash
	}

	if len(roles) == 0 {
		roles = append(roles, dtypes.UserRoleUser)
	}

	createUser := func(ctx context.Context) error {
		var err error
		user, err = a.createUser(ctx, req.Login, req.Email, passwordHash, roles)
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

func (a *Auth) createUser(
	ctx context.Context,
	login string,
	email string,
	passwordHash *string,
	roles []dtypes.UserRole,
) (*dauth.User, error) {
	userID, err := a.user.CreateUser(ctx, login, email, passwordHash, roles)
	if err != nil {
		return nil, err
	}

	user, err := a.user.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
