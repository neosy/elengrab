package authz

import (
	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func IsAnonymousByUserID(userID uuid.UUID) bool {
	return userID == dauth.AnonymousUserID()
}

func IsAnonymous(roles dtypes.UserRoleIDs) bool {
	return len(roles) == 0
}

func IsAdmin(roles dtypes.UserRoleIDs) bool {
	return roles.HasRoleID(dtypes.UserRoleAdmin)
}

func IsGuest(roles dtypes.UserRoleIDs) bool {
	if IsAdmin(roles) || IsUser(roles) {
		return false
	}
	return roles.HasRoleID(dtypes.UserRoleGuest)
}

func IsUser(roles dtypes.UserRoleIDs) bool {
	return !IsAnonymous(roles) && !IsGuest(roles)
}
