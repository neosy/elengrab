package useruc

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/pkg/errorx"
	"github.com/neosy/elengrab/pkg/errorx/exceptionx"
)

type UserOption func(*dauth.User)

func (uc *User) WithGuest(isGuest bool) UserOption {
	return func(u *dauth.User) {
		u.IsGuest = isGuest
	}
}

func (uc *User) Create(ctx context.Context, user *dauth.User, opts ...UserOption) (uuid.UUID, error) {
	if user == nil {
		uc.logger.Warn("Nil pointer in function")
		return uuid.Nil, errors.New("function parameter is a null pointer")
	}

	if user.UserID == uuid.Nil {
		user.UserID = uuid.New()
	}

	for _, opt := range opts {
		opt(user)
	}

	err := uc.userRep.Insert(ctx, user)
	if err != nil {
		uc.logger.Warn(
			"Failed to insert record into repository",
			"error", err,
		)
		return uuid.Nil, err
	}

	return user.UserID, nil
}

func (uc *User) CreateGuest(ctx context.Context) (uuid.UUID, error) {
	login, err := uc.genLogin()
	if err != nil {
		return uuid.Nil, errorx.NewByErr(err, exceptionx.ERROR)
	}

	user := &dauth.User{
		Login:    login,
		Email:    nil,
		IsActive: true,
	}

	return uc.Create(ctx, user, uc.WithGuest(true))
}

func (u *User) genLogin() (string, error) {
	b := make([]byte, 2)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	ts := strconv.FormatInt(time.Now().Unix(), 36)

	return fmt.Sprintf("user-%s%s", hex.EncodeToString(b), ts), nil
}
