package dauth

import (
	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func ResolveUserType(userID uuid.UUID, roleIDs dtypes.UserRoleIDs) dtypes.UserType {
	if userID == AnonymousUserID() {
		return dtypes.UserTypeAnonymous
	}
	return dtypes.ResolveUserTypeByRoles(roleIDs)
}
