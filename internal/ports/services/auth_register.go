package pservices

import dtypes "github.com/neosy/elengrab/internal/domain/types"

type AuthRegisterAdminRequest struct {
	Login    string
	Email    string
	Password string
}

type AuthRegisterUserRequest struct {
	Login    string
	Email    string
	Password string
	RoleIDs  dtypes.UserRoleIDs
}
