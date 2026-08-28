package authsession

import "context"

func (uc *UserSession) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.userSessionRepo().Tx(ctx, fn)
}
