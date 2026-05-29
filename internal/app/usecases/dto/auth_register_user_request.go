package dto

import dtypes "github.com/neosy/elengrab/internal/domain/types"

type RegisterUserRequest struct {
	Login    string
	Email    string
	Password string
	RoleIDs  dtypes.UserRoleIDs
}
