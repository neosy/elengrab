package admin

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/admin/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

type Admin interface {
	SetUserRoles(ctx context.Context, req dto.SetUserRolesRequest) error

	GetUserInfo(ctx context.Context, userID uuid.UUID) (*dto.UserInfoResponse, error)
	GetAllUsersWithoutGuest(ctx context.Context) ([]*dauth.User, error)
}
