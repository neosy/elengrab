package authuser

import (
	"context"

	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

// FindByUserID
func (u *User) FindByUserID(ctx context.Context, userID uuid.UUID) (*dauth.User, error) {
	if userID == uuid.Nil {
		return nil, nil
	}

	repo := u.userRepo().WithoutDeleted()

	user, err := repo.FindByUserID(ctx, userID)
	if err != nil {
		u.logger.Warn("Failed get user", "error", err)
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	return user, nil
}

func (u *User) FindByLogin(ctx context.Context, login dtypes.Login) (*dauth.User, error) {
	if login == "" {
		return nil, nil
	}

	repo := u.userRepo().WithoutDeleted()

	user, err := repo.FindByLogin(ctx, login)
	if err != nil {
		u.logger.Warn("Failed get user", "error", err)
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	return user, nil
}

func (u *User) FindByEmail(ctx context.Context, email string) (*dauth.User, error) {
	if email == "" {
		return nil, nil
	}

	repo := u.userRepo().WithoutDeleted()

	user, err := repo.FindByEmail(ctx, email)
	if err != nil {
		u.logger.Warn("Failed get user", "error", err)
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	return user, nil
}

// GetByUserID
// Record MUST exist — otherwise NOT_FOUND
func (u *User) GetByUserID(ctx context.Context, userID uuid.UUID) (*dauth.User, error) {
	user, err := u.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		u.logger.Debug("User not found", "userID", userID)
		return nil, errorx.New("user not found", exceptionx.NOT_FOUND)
	}

	return user, nil
}

func (u *User) GetByLogin(ctx context.Context, login dtypes.Login) (*dauth.User, error) {
	user, err := u.FindByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	if user == nil {
		u.logger.Debug("User not found", "login", login)
		return nil, errorx.New("user not found", exceptionx.NOT_FOUND)
	}

	return user, nil
}

func (u *User) GetByEmail(ctx context.Context, email string) (*dauth.User, error) {
	user, err := u.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		u.logger.Debug("User not found", "email", email)
		return nil, errorx.New("user not found", exceptionx.NOT_FOUND)
	}

	return user, nil
}

func (u *User) ExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error) {
	repo := u.userRepo().WithoutDeleted()

	exists, err := repo.ExistsByUserID(ctx, userID)
	if err != nil {
		u.logger.Warn("Failed to check if user exists", "userID", userID, "error", err)
	}

	return exists, nil
}

func (u *User) ExistsByLogin(ctx context.Context, login dtypes.Login) (bool, error) {
	repo := u.userRepo().WithoutDeleted()

	exists, err := repo.ExistsByLogin(ctx, login)
	if err != nil {
		u.logger.Warn("Failed to check if user exists", "login", login, "error", err)
	}

	return exists, nil
}

func (u *User) GetAllUsers(ctx context.Context) ([]*dauth.User, error) {
	var users []*dauth.User

	rep := u.userRepo().WithoutDeleted()

	err := rep.IterateGetAll(
		ctx,
		func(u *dauth.User) error {
			u.PasswordHash = nil
			users = append(users, u)
			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *User) GetAllUsersWithoutGuest(ctx context.Context) ([]*dauth.User, error) {
	var users []*dauth.User

	rep := u.userRepo().WithoutDeleted().WithoutGuest()

	err := rep.IterateGetAll(
		ctx,
		func(u *dauth.User) error {
			u.PasswordHash = nil
			users = append(users, u)
			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return users, nil
}
