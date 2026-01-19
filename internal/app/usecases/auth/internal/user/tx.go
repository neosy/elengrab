package useruc

import "context"

func (uc *User) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.userRep.Tx(ctx, fn)
}
