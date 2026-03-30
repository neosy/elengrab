package authuserrole

import (
	"context"

	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

// Find
func (uc *UserRole) Find(ctx context.Context, userID uuid.UUID, roleID string) (*dauth.UserRole, error) {
	if userID == uuid.Nil || roleID == "" {
		return nil, nil
	}

	userRole, err := uc.userRoleRep.Find(ctx, userID, roleID)
	if err != nil {
		uc.logger.Warn("Failed get userRole", "error", err)
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	return userRole, nil
}

// Get
// Record MUST exist — otherwise NOT_FOUND
func (uc *UserRole) Get(ctx context.Context, userID uuid.UUID, roleID string) (*dauth.UserRole, error) {
	userRole, err := uc.Find(ctx, userID, roleID)
	if err != nil {
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	if userRole == nil {
		uc.logger.Debug("UserRole not found", "userID", userID, "roleID", roleID)
		return nil, errorx.New("userRole not found", exceptionx.NOT_FOUND)
	}

	return userRole, nil
}

// Exists
func (uc *UserRole) Exists(ctx context.Context, userID uuid.UUID, roleID string) (bool, error) {
	exists, err := uc.userRoleRep.Exists(ctx, userID, roleID)
	if err != nil {
		uc.logger.Warn(
			"Failed to check if userRole exists",
			"userID", userID,
			"roleID", roleID,
			"error", err,
		)
	}

	return exists, nil
}
