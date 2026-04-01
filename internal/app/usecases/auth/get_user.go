package auth

import (
	"context"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (a *Auth) ExistsUserByLogin(ctx context.Context, login string) (bool, error) {
	return a.user.ExistsByLogin(ctx, dtypes.NewLogin(login))
}
