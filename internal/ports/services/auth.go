package pservices

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
)

type AuthMiddleware interface {
	RegisterGuest(ctx context.Context) (*dto.AuthUserResponse, error)

	ValidateSession(ctx context.Context, token string) (*dto.AuthUserResponse, error)
	RefreshSession(ctx context.Context, token string) (*dto.AuthUserResponse, error)
}

type AuthService interface {
	SoftDeleteUser(ctx context.Context, userID uuid.UUID) error

	RegisterUser(ctx context.Context, req *dto.RegisterUserRequest) (*dto.AuthUserResponse, error)
	RegisterGuest(ctx context.Context) (*dto.AuthUserResponse, error)
	RegisterAdmin(ctx context.Context, req *dto.RegisterAdminRequest) (*dto.AuthUserResponse, error)

	AuthenticateUser(ctx context.Context, req *dto.AuthUserRequest) (*dto.AuthUserResponse, error)

	ExistsUserByLogin(ctx context.Context, login string) (bool, error)
}
