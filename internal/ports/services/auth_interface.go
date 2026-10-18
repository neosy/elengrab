package pservices

import (
	"context"

	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

type AuthMiddleware interface {
	RegisterGuest(ctx context.Context) (*AuthUserResponse, error)

	ValidateSession(ctx context.Context, token string) (*AuthUserResponse, error)
	RefreshSession(ctx context.Context, token string) (*AuthUserResponse, error)
}

type AuthService interface {
	SoftDeleteUser(ctx context.Context, userID uuid.UUID) error

	RegisterUser(ctx context.Context, req *AuthRegisterUserRequest) (*AuthUserResponse, error)
	RegisterGuest(ctx context.Context) (*AuthUserResponse, error)
	RegisterAdmin(ctx context.Context, req *AuthRegisterAdminRequest) (*AuthUserResponse, error)

	AuthenticateUser(ctx context.Context, req *AuthUserRequest) (*AuthUserResponse, error)

	FindByUserID(ctx context.Context, userID uuid.UUID) (*dauth.User, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*dauth.User, error)
	FindByLogin(ctx context.Context, login string) (*dauth.User, error)
	GetByLogin(ctx context.Context, login string) (*dauth.User, error)
	ExistsUserByLogin(ctx context.Context, login string) (bool, error)
}

type AdminPanelService interface {
	SetUserRoles(ctx context.Context, userID uuid.UUID, roleIDs []string) error

	GetAllUsers(ctx context.Context) ([]*dauth.User, error)
	GetAllUsersWithoutGuest(ctx context.Context) ([]*dauth.User, error)

	FindByUserID(ctx context.Context, userID uuid.UUID) (*dauth.User, error)

	GetAllRoles(ctx context.Context) ([]*dauth.Role, error)
	GetAllRolesWithoutGuest(ctx context.Context) ([]*dauth.Role, error)
}
