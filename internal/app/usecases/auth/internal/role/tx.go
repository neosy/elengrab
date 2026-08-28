package authrole

import "context"

func (uc *Role) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.roleRepo().Tx(ctx, fn)
}
