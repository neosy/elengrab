package pservices

import (
	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type AuthUserRequest struct {
	Login    string
	Email    string
	Password string
}

type AuthUserResponse struct {
	UserID  uuid.UUID
	Login   string
	Email   string
	RoleIDs []string
	Token   *dtypes.AuthToken
}

func (u *AuthUserResponse) UserType() dtypes.UserType {
	roleIDs, _ := dtypes.ParseUserRoleIDs(u.RoleIDs)
	return dauth.ResolveUserType(u.UserID, roleIDs)
}
