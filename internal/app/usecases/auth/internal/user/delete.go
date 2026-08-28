package authuser

import (
	"context"

	"github.com/google/uuid"
)

func (u *User) SoftDelete(ctx context.Context, userID uuid.UUID) error {
	err := u.userRepo().Delete(ctx, userID, true)
	if err != nil {
		u.logger.Warn("Failed soft delete user", "error", err)
		return err
	}
	return nil
}

func (u *User) HardDelete(ctx context.Context, userID uuid.UUID) error {
	err := u.userRepo().Delete(ctx, userID, false)
	if err != nil {
		u.logger.Warn("Failed delete user", "error", err)
		return err
	}
	return nil
}
