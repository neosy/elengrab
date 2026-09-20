package authz

import (
	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func IsAnonymousByUserID(userID uuid.UUID) bool {
	userType := dauth.ResolveUserType(userID, nil)
	return userType == dtypes.UserTypeAnonymous
}

func IsAnonymous(roleIDs dtypes.UserRoleIDs) bool {
	userType := dtypes.ResolveUserTypeByRoles(roleIDs)
	return userType == dtypes.UserTypeAnonymous
}

func IsAdmin(roleIDs dtypes.UserRoleIDs) bool {
	userType := dtypes.ResolveUserTypeByRoles(roleIDs)
	return userType == dtypes.UserTypeAdmin
}

func IsGuest(roleIDs dtypes.UserRoleIDs) bool {
	userType := dtypes.ResolveUserTypeByRoles(roleIDs)
	return userType == dtypes.UserTypeGuest
}

func IsUser(roleIDs dtypes.UserRoleIDs) bool {
	userType := dtypes.ResolveUserTypeByRoles(roleIDs)
	return userType == dtypes.UserTypeUser
}
