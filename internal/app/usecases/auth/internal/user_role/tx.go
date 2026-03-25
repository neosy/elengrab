package authuserrole

import "context"

func (uc *UserRole) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.userRoleRep.Tx(ctx, fn)
}
