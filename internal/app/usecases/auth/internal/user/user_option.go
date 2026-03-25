package authuser

import dtypes "github.com/neosy/elengrab/internal/domain/types"

type UserOption func(*[]dtypes.UserRole)

func GuestRoleOption() UserOption {
	return RolesOption(dtypes.UserRoleGuest)
}

func RolesOption(roles ...dtypes.UserRole) UserOption {
	return func(rs *[]dtypes.UserRole) {
		*rs = append(*rs, roles...)
	}
}
