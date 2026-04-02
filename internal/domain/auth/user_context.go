package dauth

import (
	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type UserContext struct {
	UserID uuid.UUID
	Login  string
	Email  string
	Roles  []dtypes.UserRole
}

func UserContextAnonymous() *UserContext {
	return &UserContext{
		UserID: uuid.Nil,
		Roles:  []dtypes.UserRole{dtypes.UserRoleGuest},
	}
}

func (u *UserContext) UserType() dtypes.UserType {
	if u.UserID == uuid.Nil {
		return dtypes.UserTypeAnonymous
	}

	var uType dtypes.UserType = dtypes.UserTypeAnonymous

	for _, r := range u.Roles {
		switch r {
		case dtypes.UserRoleAdmin:
			return dtypes.UserTypeAdmin
		case dtypes.UserRoleGuest:
			if uType < dtypes.UserTypeGuest {
				uType = dtypes.UserTypeGuest
			}
		default:
			uType = dtypes.UserTypeUser
		}
	}

	return uType
}

func (u *UserContext) IsGuest() bool {
	return u.UserType() == dtypes.UserTypeGuest
}

func (u *UserContext) IsUserTypeAnonymous() bool {
	return u.UserType() == dtypes.UserTypeAnonymous
}

func (u *UserContext) IsUser() bool {
	return u.UserType() == dtypes.UserTypeUser
}

func (u *UserContext) IsAdmin() bool {
	return u.UserType() == dtypes.UserTypeAdmin
}
