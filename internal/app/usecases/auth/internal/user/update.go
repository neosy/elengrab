package authuser

import (
	"context"

	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

func (u *User) Update(ctx context.Context, user *dauth.User) error {
	err := u.userRep.Update(ctx, user)
	if err != nil {
		u.logger.Warn("Update record error", "error", err)
		return err
	}
	return nil
}
