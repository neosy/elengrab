package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

type Auth interface {
	RegisterUser(ctx context.Context, req *dto.RegisterUserRequest) (*dto.AuthUserResponse, error)
	RegisterGuest(ctx context.Context) (*dto.AuthUserResponse, error)
	RegisterAdmin(ctx context.Context, req *dto.RegisterAdminRequest) (*dto.AuthUserResponse, error)

	SoftDeleteUser(ctx context.Context, userID uuid.UUID) error

	SetUserRoles(ctx context.Context, userID uuid.UUID, roleIDs []string) error

	ValidateSession(ctx context.Context, token string) (*dto.AuthUserResponse, error)
	RefreshSession(ctx context.Context, token string) (*dto.AuthUserResponse, error)

	AuthenticateUser(ctx context.Context, req *dto.AuthUserRequest) (*dto.AuthUserResponse, error)

	FindByUserID(ctx context.Context, userID uuid.UUID) (*dauth.User, error)
	ExistsUserByLogin(ctx context.Context, login string) (bool, error)
	GetAllUsers(ctx context.Context) ([]*dauth.User, error)
	GetAllUsersWithoutGuest(ctx context.Context) ([]*dauth.User, error)

	GetAllRoles(ctx context.Context) ([]*dauth.Role, error)
	GetAllRolesWithoutGuest(ctx context.Context) ([]*dauth.Role, error)
}
