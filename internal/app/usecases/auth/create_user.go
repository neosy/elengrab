package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/pkg/errorx"
	"github.com/neosy/elengrab/pkg/errorx/exceptionx"
)

func (u *Auth) CreateUser(ctx context.Context) (*dauth.User, error) {
	login, err := u.genLogin()
	if err != nil {
		return nil, errorx.NewByErr(err, exceptionx.ERROR)
	}

	user := &dauth.User{
		Login:    login,
		Email:    nil,
		IsActive: true,
	}

	userID, err := u.user.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	user, err = u.user.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *Auth) genLogin() (string, error) {
	b := make([]byte, 2)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	ts := strconv.FormatInt(time.Now().Unix(), 36)

	return fmt.Sprintf("user-%s%s", hex.EncodeToString(b), ts), nil
}
