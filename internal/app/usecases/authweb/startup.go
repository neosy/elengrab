package authweb

import (
	"context"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
)

func (a *AuthWeb) Startup(ctx context.Context) error {
	err := a.registerDefaultAdmin(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthWeb) registerDefaultAdmin(ctx context.Context) error {
	if a.defaultAdminLogin == "" || a.defaultAdminPassword == "" {
		return nil
	}

	exists, err := a.auth.ExistsUserByLogin(ctx, a.defaultAdminLogin)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	req := &dto.RegisterAdminRequest{
		Login:    a.defaultAdminLogin,
		Password: a.defaultAdminPassword,
	}

	_, err = a.auth.RegisterAdmin(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
