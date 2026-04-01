package auth

import (
	"context"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (a *Auth) AuthenticateUser(
	ctx context.Context,
	req *dto.AuthUserRequest,
) (*dto.AuthUserResponse, error) {
	if req == nil {
		return nil, errorx.New("function parameter is nil", exceptionx.ERROR)
	}

	var (
		user    *dauth.User
		session *dauth.UserSession
		err     error
	)

	login := dtypes.NewLogin(req.Login)

	if login == "" && req.Email == "" {
		return nil, errorx.New("login or email is required", exceptionx.VALIDATE)
	}

	if login != "" {
		user, err = a.user.FindByLogin(ctx, login)
		if err != nil {
			return nil, err
		}
	}

	if user == nil && req.Email != "" {
		user, err = a.user.FindByEmail(ctx, req.Email)
		if err != nil {
			return nil, err
		}
	}

	if user == nil {
		a.logger.Debug("User not found", "login", login, "email", req.Email)
		return nil, ErrUserNotFound
	}

	if !user.IsActive {
		a.logger.Info("User is not active", "userID", user.UserID)
		return nil, ErrUserIsNotActive
	}

	if user.DeletedAt != nil {
		a.logger.Info("User deleted", "userID", user.UserID)
		return nil, ErrUserDeleted
	}

	if user.PasswordHash == nil {
		a.logger.Info("The user is not available because the password is not set", "userID", user.UserID)
		return nil, errorx.New("the user is not available because the password is not set.", exceptionx.UNAUTHORIZED)
	}

	if err := a.checkPassword(*user.PasswordHash, req.Password); err != nil {
		a.logger.Info("check password failed", "error", err)
		return nil, errorx.NewFromError(err, exceptionx.UNAUTHORIZED)
	}

	session, err = a.createSession(ctx, user.UserID)
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
