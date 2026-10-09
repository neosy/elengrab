package dto

import (
	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type AuthUserResponse struct {
	UserID  uuid.UUID
	Login   string
	Email   string
	RoleIDs []string
	Token   *AuthToken
}

func (u *AuthUserResponse) UserType() dtypes.UserType {
	roleIDs, _ := dtypes.ParseUserRoleIDs(u.RoleIDs)
	return dauth.ResolveUserType(u.UserID, roleIDs)
}
