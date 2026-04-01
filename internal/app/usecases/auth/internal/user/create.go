package authuser

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (u *User) Create(ctx context.Context, user *dauth.User, opts ...UserOption) (uuid.UUID, error) {
	if user == nil {
		u.logger.Warn("Nil pointer in function")
		return uuid.Nil, errors.New("function parameter is a null pointer")
	}

	if user.UserID == uuid.Nil {
		user.UserID = uuid.New()
	}

	var roles []dtypes.UserRole

	for _, opt := range opts {
		opt(&roles)
	}

	err := u.userRep.Tx(
		ctx,
		func(ctx context.Context) error {
			err := u.userRep.Insert(ctx, user)
			if err != nil {
				return err
			}

			for _, role := range roles {
				err := u.userRole.Create(
					ctx,
					&dauth.UserRole{
						UserID: user.UserID,
						RoleID: role.String(),
					},
				)
				return err
			}

			return nil
		},
	)

	if err != nil {
		u.logger.Warn(
			"Failed to insert record into repository",
			"error", err,
		)
		return uuid.Nil, err
	}

	return user.UserID, nil
}

func (u *User) CreateGuest(ctx context.Context) (uuid.UUID, error) {
	return u.CreateUser(ctx, "", "", nil, nil)
}

func (u *User) CreateAdmin(ctx context.Context, login dtypes.Login, email string, passwordHash *string) (uuid.UUID, error) {
	if login == "" {
		login = dtypes.UserRoleAdmin.Login()
	}
	return u.CreateUser(ctx, login, email, passwordHash, []dtypes.UserRole{dtypes.UserRoleAdmin})
}

func (u *User) CreateUser(
	ctx context.Context,
	login dtypes.Login,
	email string,
	passwordHash *string,
	roles []dtypes.UserRole,
) (uuid.UUID, error) {
	if login.String() == "" {
		var err error
		l, err := u.genLogin()
		if err != nil {
			return uuid.Nil, errorx.NewFromError(err, exceptionx.ERROR)
		}
		login = l
	}

	var emailPtr *string
	if email != "" {
		emailPtr = &email
	}

	user := &dauth.User{
		Login:        login,
		Email:        emailPtr,
		PasswordHash: passwordHash,
		IsActive:     true,
	}

	var rolesOpt UserOption
	if len(roles) > 0 {
		rolesOpt = RolesOption(roles...)
	} else {
		rolesOpt = GuestRoleOption()
	}

	return u.Create(ctx, user, rolesOpt)
}

func (u *User) genLogin() (dtypes.Login, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	rndPart := strconv.FormatUint(uint64(binary.BigEndian.Uint32(b)), 36)

	ts := strconv.FormatInt(time.Now().UTC().Unix(), 36)

	return dtypes.NewLogin(fmt.Sprintf("u-%s%s", ts, rndPart)), nil
}
