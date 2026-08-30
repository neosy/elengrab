package auth

import (
	"context"
	"strings"

	autherr "github.com/neosy/elengrab/internal/app/usecases/auth/errors"
	"github.com/neosy/elengrab/internal/app/usecases/auth/internal/consts"
	idto "github.com/neosy/elengrab/internal/app/usecases/auth/internal/dto"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (a *auth) RegisterUser(
	ctx context.Context,
	req *dto.RegisterUserRequest,
) (*dto.AuthUserResponse, error) {
	if req == nil {
		return nil, autherr.ErrFunctionNilParameter
	}

	var (
		user         *dauth.User
		session      *dauth.UserSession
		passwordHash *string
		roles        []string
	)

	login := strings.TrimSpace(req.Login)
	email := strings.TrimSpace(req.Email)
	password := strings.TrimSpace(req.Password)

	if password != "" {
		hash, err := a.hashPassword(password)
		if err != nil {
			return nil, err
		}
		passwordHash = &hash
	}

	if len(roles) == 0 {
		roles = append(roles, consts.UserRole)
	}

	createUser := func(ctx context.Context) error {
		var err error
		createReq := &idto.CreateUserRequest{
			Login:        dtypes.NewLogin(login),
			Email:        email,
			PasswordHash: passwordHash,
			Roles:        roles,
		}
		user, err = a.createUser(ctx, createReq)
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

func (a *auth) createUser(ctx context.Context, req *idto.CreateUserRequest) (*dauth.User, error) {
	userID, err := a.user.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	user, err := a.user.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
