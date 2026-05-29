package dto

import (
	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type SetUserRolesRequest struct {
	UserID  uuid.UUID
	RoleIDs dtypes.UserRoleIDs
}
