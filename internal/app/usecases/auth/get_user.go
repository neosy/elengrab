package auth

import "context"

func (a *Auth) ExistsUserByLogin(ctx context.Context, login string) (bool, error) {
	return a.user.ExistsByLogin(ctx, login)
}
