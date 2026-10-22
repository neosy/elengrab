package auth

import (
	"context"

	"uuid"
)

func (a *auth) SoftDeleteUser(ctx context.Context, userID uuid.UUID) error {
	return a.user.SoftDelete(ctx, userID)
}
