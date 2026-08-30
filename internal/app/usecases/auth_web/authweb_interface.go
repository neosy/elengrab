package authweb

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
)

type AuthWeb interface {
	Startup(ctx context.Context) error

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
