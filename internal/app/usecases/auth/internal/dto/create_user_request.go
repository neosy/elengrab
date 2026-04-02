package idto

import dtypes "github.com/neosy/elengrab/internal/domain/types"

type CreateUserRequest struct {
	Login        dtypes.Login
	Email        string
	PasswordHash *string
	Roles        []dtypes.UserRole
}
