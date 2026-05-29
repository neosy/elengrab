package authuser

import (
	"github.com/neosy/elengrab/internal/app/usecases/auth/internal/consts"
)

type UserOption func(*[]string)

func GuestRoleOption() UserOption {
	return RolesOption(consts.GuestRole)
}

func AdminRoleOption() UserOption {
	return RolesOption(consts.AdminRole)
}

func RolesOption(roles ...string) UserOption {
	return func(rs *[]string) {
		*rs = append(*rs, roles...)
	}
}
