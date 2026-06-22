package dauth

import (
	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

var AnonymousUserID = func() uuid.UUID { return uuid.Nil }

type UserContext struct {
	UserID        uuid.UUID
	AnonSessionID uuid.UUID

	Login        string
	Email        string
	RoleIDs      dtypes.UserRoleIDs
	GuestCreated bool
}

func UserContextAnonymous(anonSessionID uuid.UUID) UserContext {
	return UserContext{
		UserID:        AnonymousUserID(),
		AnonSessionID: anonSessionID,
		RoleIDs:       nil,
	}
}

func (u *UserContext) UserType() dtypes.UserType {
	if u == nil || u.UserID == AnonymousUserID() {
		return dtypes.UserTypeAnonymous
	}

	var uType dtypes.UserType = dtypes.UserTypeAnonymous

	for _, r := range u.RoleIDs {
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

func (u *UserContext) EventKey() dtypes.EventKey {
	if u.UserID != uuid.Nil {
		return dtypes.NewEventKeyUserID(u.UserID)
	}
	return dtypes.NewEventKeySessionID(u.AnonSessionID)
}

func (u *UserContext) IsGuest() bool {
	return u.UserType() == dtypes.UserTypeGuest
}

func (u *UserContext) IsAnonymous() bool {
	return u.UserType() == dtypes.UserTypeAnonymous
}

func (u *UserContext) IsUser() bool {
	return u.UserType() == dtypes.UserTypeUser
}

func (u *UserContext) IsAdmin() bool {
	return u.UserType() == dtypes.UserTypeAdmin
}
