package authuser

import "context"

func (u *User) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return u.userRep.Tx(ctx, fn)
}
