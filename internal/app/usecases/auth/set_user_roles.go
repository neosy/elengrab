package auth

import (
	"context"
	"slices"

	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (a *Auth) SetUserRoles(ctx context.Context, userID uuid.UUID, roleIDs []string) error {
	if len(roleIDs) == 0 {
		return errorx.NewMessage("roles are not defined", exceptionx.WRONG_DATA)
	}

	user, err := a.user.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	currentRoleIDs := user.RoleIDs

	if len(currentRoleIDs) == len(roleIDs) {
		a := slices.Clone(currentRoleIDs)
		b := slices.Clone(roleIDs)

		slices.Sort(a)
		slices.Sort(b)

		if slices.Equal(a, b) {
			return errorx.NewMessage("user roles are unchanged", exceptionx.WRONG_DATA)
		}
	}

	setRoles := func(ctx context.Context) error {
		for _, roleID := range roleIDs {
			if slices.Contains(currentRoleIDs, roleID) {
				continue
			}
			userRole := &dauth.UserRole{
				UserID: userID,
				RoleID: roleID,
			}
			err := a.userRole.Create(ctx, userRole)
			if err != nil {
				return err
			}
		}

		for _, roleID := range currentRoleIDs {
			if !slices.Contains(roleIDs, roleID) {
				err := a.userRole.Delete(ctx, userID, roleID)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}

	return a.userRole.Tx(ctx, setRoles)
}
