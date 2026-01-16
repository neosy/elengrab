package useruc

import (
	"context"

	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

func (uc *User) Update(ctx context.Context, user *dauth.User) error {
	err := uc.userRep.Update(ctx, user)
	if err != nil {
		uc.logger.Warn("Update record error", "error", err)
		return err
	}
	return nil
}
