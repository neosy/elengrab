package useruc

import (
	"context"

	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/pkg/errorx"
	"github.com/neosy/elengrab/pkg/errorx/exceptionx"
)

// FindByUserID
func (uc *User) FindByUserID(ctx context.Context, userID uuid.UUID) (*dauth.User, error) {
	if userID == uuid.Nil {
		return nil, nil
	}

	user, err := uc.userRep.FindByUserID(ctx, userID)
	if err != nil {
		uc.logger.Warn("Failed get user", "error", err)
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	return user, nil
}

// GetByUserID
// Record MUST exist — otherwise NOT_FOUND
func (uc *User) GetByUserID(ctx context.Context, userID uuid.UUID) (*dauth.User, error) {
	user, err := uc.FindByUserID(ctx, userID)
	if err != nil {
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	if user == nil {
		uc.logger.Warn("User not found", "userID", userID)
		return nil, errorx.New("user not found", exceptionx.NOT_FOUND)
	}

	return user, nil
}

func (uc *User) ExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error) {
	exists, err := uc.userRep.ExistsByUserID(ctx, userID)
	if err != nil {
		uc.logger.Warn("Failed to check if user exists", "userID", userID, "error", err)
	}

	return exists, nil
}
