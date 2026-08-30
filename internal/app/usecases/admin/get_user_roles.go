package admin

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/admin/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (a *admin) GetUserInfo(ctx context.Context, userID uuid.UUID) (*dto.UserInfoResponse, error) {
	user, err := a.adminPanel.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		a.logger.Debug("User not found", "userID", userID)
		return nil, errorx.New("user not found", exceptionx.NOT_FOUND)
	}

	roles, err := a.adminPanel.GetAllRolesWithoutGuest(ctx)
	if err != nil {
		return nil, err
	}

	rolesByID := roleIDsToMap(user.RoleIDs)

	rolesWithAssignment := make([]dto.RoleWithAssignment, 0, len(roles))
	for _, role := range roles {
		_, assigned := rolesByID[dtypes.UserRoleID(role.RoleID)]

		rolesWithAssignment = append(
			rolesWithAssignment,
			dto.RoleWithAssignment{
				RoleID:   role.RoleID,
				Name:     role.Name,
				Assigned: assigned,
			},
		)
	}

	return &dto.UserInfoResponse{
		User:                *user,
		RolesWithAssignment: rolesWithAssignment,
	}, nil
}

func roleIDsToMap(roleIDs []string) map[dtypes.UserRoleID]struct{} {
	rolesByID := make(map[dtypes.UserRoleID]struct{}, len(roleIDs))
	for _, roleID := range roleIDs {
		rolesByID[dtypes.UserRoleID(roleID)] = struct{}{}
	}
	return rolesByID
}
