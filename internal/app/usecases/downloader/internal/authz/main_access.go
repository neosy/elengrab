package authz

import (
	iconfig "github.com/neosy/elengrab/internal/config"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// HasViewAllAccess returns true if the roles allow viewing all (global access)
func (a *Authorization) HasViewAllAccess(roles dtypes.UserRoleIDs) bool {
	if a.appMode == dtypes.AppModePublic {
		return true
	}
	return roles.HasAnyRoleID(dtypes.UserRoleAdmin, dtypes.UserRoleViewerAll)
}

// HasCreateAccess reports whether the user is allowed to create new records.
func (a *Authorization) HasCreateAccess(roles dtypes.UserRoleIDs) bool {
	if iconfig.DemoMode() {
		return false
	}

	if a.appMode == dtypes.AppModePublic {
		return true
	}
	if a.appMode == dtypes.AppModePublicReadonly || a.appMode == dtypes.AppModeGuest {
		return !IsAnonymous(roles)
	}
	return roles.HasAnyRoleID(dtypes.UserRoleAdmin, dtypes.UserRoleViewerAll)
}

// HasWriteAllAccess reports whether the user has full write access,
// including creating, updating, and deleting records.
func (a *Authorization) HasWriteAllAccess(roles dtypes.UserRoleIDs) bool {
	if iconfig.DemoMode() {
		return false
	}

	if a.appMode == dtypes.AppModePublic {
		return true
	}
	return roles.HasAnyRoleID(dtypes.UserRoleAdmin, dtypes.UserRoleViewerAll)
}

func (a *Authorization) HasPublicViewAccess(roles dtypes.UserRoleIDs) bool {
	return true
}
