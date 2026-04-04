package authweb

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/valyala/fasthttp"
)

func (a *AuthWeb) SoftDeleteUser(ctx context.Context, userID uuid.UUID) error {
	if userID == uuid.Nil {
		a.logger.Warn("UserID is Nil")
		return errorx.NewHTTP("invalid userID", fasthttp.StatusBadRequest)
	}
	if err := a.auth.SoftDeleteUser(ctx, userID); err != nil {
		a.logger.Warn("Error deleting user", "userID", userID, "error", err)
		return err
	}
	return nil
}
