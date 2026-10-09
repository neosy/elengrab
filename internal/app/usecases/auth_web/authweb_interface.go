package authweb

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
)

type AuthWeb interface {
	Startup(ctx context.Context) error

	FindByUserID(ctx context.Context, userID uuid.UUID) (*dto.AuthUserResponse, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*dto.AuthUserResponse, error)
	FindByLogin(ctx context.Context, login string) (*dto.AuthUserResponse, error)
	GetByLogin(ctx context.Context, login string) (*dto.AuthUserResponse, error)

	RegisterUser(
		ctx context.Context,
		req *dto.RegisterUserRequest,
	) (*dto.AuthUserResponse, error)

	LoginUser(
		ctx context.Context,
		req *dto.AuthUserRequest,
	) (*dto.AuthUserResponse, error)

	SoftDeleteUser(ctx context.Context, userID uuid.UUID) error
}
