package dto

import (
	"uuid"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type SetUserRolesRequest struct {
	UserID  uuid.UUID
	RoleIDs dtypes.UserRoleIDs
}
