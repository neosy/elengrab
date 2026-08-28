package authrole

import (
	"context"

	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

// Find
func (uc *Role) Find(ctx context.Context, roleID string) (*dauth.Role, error) {
	if roleID == "" {
		return nil, nil
	}

	role, err := uc.roleRepo().Find(ctx, roleID)
	if err != nil {
		uc.logger.Warn("Failed get role", "error", err)
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	return role, nil
}

// Get
// Record MUST exist — otherwise NOT_FOUND
func (uc *Role) Get(ctx context.Context, roleID string) (*dauth.Role, error) {
	role, err := uc.Find(ctx, roleID)
	if err != nil {
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	if role == nil {
		uc.logger.Debug("Role not found", "roleID", roleID)
		return nil, errorx.New("role not found", exceptionx.NOT_FOUND)
	}

	return role, nil
}

// Exists
func (uc *Role) Exists(ctx context.Context, roleID string) (bool, error) {
	exists, err := uc.roleRepo().Exists(ctx, roleID)
	if err != nil {
		uc.logger.Warn("Failed to check if role exists", "roleID", roleID, "error", err)
	}

	return exists, nil
}

func (uc *Role) GetAll(ctx context.Context) ([]*dauth.Role, error) {
	return uc.roleRepo().GetAll(ctx)
}

func (uc *Role) GetAllWithoutGuest(ctx context.Context) ([]*dauth.Role, error) {
	return uc.roleRepo().WithoutGuest().GetAll(ctx)
}
