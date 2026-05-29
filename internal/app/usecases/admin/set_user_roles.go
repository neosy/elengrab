package admin

import (
	"context"

	"github.com/neosy/elengrab/internal/app/usecases/admin/dto"
)

func (a *Admin) SetUserRoles(ctx context.Context, req dto.SetUserRolesRequest) error {
	return a.adminPanel.SetUserRoles(ctx, req.UserID, req.RoleIDs.Strings())
}
