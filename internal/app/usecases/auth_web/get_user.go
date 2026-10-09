package authweb

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
)

func (a *authWeb) FindByUserID(ctx context.Context, userID uuid.UUID) (*dto.AuthUserResponse, error) {
	user, err := a.auth.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	userResponse := a.mappers.MapUserToUserResponse(user)

	return userResponse, nil
}

func (a *authWeb) GetByUserID(ctx context.Context, userID uuid.UUID) (*dto.AuthUserResponse, error) {
	user, err := a.auth.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	userResponse := a.mappers.MapUserToUserResponse(user)

	return userResponse, nil
}

func (a *authWeb) FindByLogin(ctx context.Context, login string) (*dto.AuthUserResponse, error) {
	user, err := a.auth.FindByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	userResponse := a.mappers.MapUserToUserResponse(user)

	return userResponse, nil
}

func (a *authWeb) GetByLogin(ctx context.Context, login string) (*dto.AuthUserResponse, error) {
	user, err := a.auth.GetByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	userResponse := a.mappers.MapUserToUserResponse(user)

	return userResponse, nil
}
