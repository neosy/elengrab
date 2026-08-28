package authrole

import (
	"context"

	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (uc *Role) Update(ctx context.Context, role *dauth.Role) error {
	err := uc.roleRepo().Update(ctx, role)
	if err != nil {
		uc.logger.Warn("Update record error", "error", err)
		return errorx.NewFromError(err, exceptionx.ERROR)
	}
	return nil
}
