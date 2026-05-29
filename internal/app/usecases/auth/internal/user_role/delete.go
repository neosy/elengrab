package authuserrole

import (
	"context"

	"github.com/google/uuid"
)

func (uc *UserRole) Delete(ctx context.Context, userID uuid.UUID, roleID string) error {
	return uc.userRoleRep.Delete(ctx, userID, roleID)
}

func (uc *UserRole) DeleteRoleIDs(ctx context.Context, userID uuid.UUID, roleIDs []string) error {
	for _, roleID := range roleIDs {
		err := uc.userRoleRep.Delete(ctx, userID, roleID)
		if err != nil {
			uc.logger.Warn(
				"Failed to delete entry from repository",
				"error", err,
			)
			return err
		}
	}
	return nil
}
