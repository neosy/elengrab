package authuser

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/neosy/elengrab/internal/app/errors"
	autherr "github.com/neosy/elengrab/internal/app/usecases/auth/errors"
	idto "github.com/neosy/elengrab/internal/app/usecases/auth/internal/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

func (u *User) Create(ctx context.Context, user *dauth.User, opts ...UserOption) (uuid.UUID, error) {
	if user == nil {
		u.logger.Warn("Nil pointer in function")
		return uuid.Nil, apperrors.ErrFuncParamNullPointer
	}

	if user.UserID == uuid.Nil {
		user.UserID = uuid.New()
	}

	var roles []string

	for _, opt := range opts {
		opt(&roles)
	}

	err := u.userRepo().Tx(
		ctx,
		func(ctx context.Context) error {
			err := u.userRepo().Insert(ctx, user)
			if err != nil {
				return err
			}

			for _, role := range roles {
				err := u.userRole.Create(
					ctx,
					&dauth.UserRole{
						UserID: user.UserID,
						RoleID: role,
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
	req := &idto.CreateUserRequest{}
	return u.CreateUser(ctx, req)
}

func (u *User) CreateAdmin(ctx context.Context, req *idto.CreateUserRequest) (uuid.UUID, error) {
	newReq := uptr.Copy(req)
	if newReq.Login == "" {
		newReq.Login = dtypes.UserRoleAdmin.Login()
	}
	return u.CreateUser(ctx, req)
}

func (u *User) CreateUser(
	ctx context.Context,
	req *idto.CreateUserRequest,
) (uuid.UUID, error) {
	if req == nil {
		return uuid.Nil, autherr.ErrFunctionNilParameter
	}

	login := req.Login
	if login == "" {
		var err error
		l, err := u.genLogin()
		if err != nil {
			return uuid.Nil, errorx.NewFromError(err, exceptionx.ERROR)
		}
		login = l
	}

	var emailPtr *string
	if req.Email != "" {
		emailPtr = &req.Email
	}

	user := &dauth.User{
		Login:        login,
		Email:        emailPtr,
		PasswordHash: req.PasswordHash,
		IsActive:     true,
	}

	var rolesOpt UserOption
	if len(req.Roles) > 0 {
		rolesOpt = RolesOption(req.Roles...)
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
