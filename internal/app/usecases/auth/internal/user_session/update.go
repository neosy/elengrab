package authsession

import (
	"context"

	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

func (uc *UserSession) Update(ctx context.Context, session *dauth.UserSession) error {
	err := uc.userSessionRepo().Update(ctx, session)
	if err != nil {
		uc.logger.Warn("Update record error", "error", err)
		return err
	}
	return nil
}
