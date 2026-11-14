package dauth

import (
	"uuid"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// ResolveUserType determines the user type based on the user ID and assigned roles.
func ResolveUserType(userID uuid.UUID, roleIDs dtypes.UserRoleIDs) dtypes.UserType {
	if userID == AnonymousUserID() {
		return dtypes.UserTypeAnonymous
	}
	return dtypes.ResolveUserTypeByRoles(roleIDs)
}
