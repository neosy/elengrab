package dauth

import (
	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	eventkey "github.com/neosy/elengrab/internal/domain/types/event_key"
)

var AnonymousUserID = func() uuid.UUID { return uuid.Nil }

type AuthContext struct {
	UserID        uuid.UUID
	AnonSessionID uuid.UUID

	Login        string
	Email        string
	RoleIDs      dtypes.UserRoleIDs
	GuestCreated bool
}

func AuthContextAnonymous(anonSessionID uuid.UUID) AuthContext {
	return AuthContext{
		UserID:        AnonymousUserID(),
		AnonSessionID: anonSessionID,
		RoleIDs:       nil,
	}
}

func (u *AuthContext) UserType() dtypes.UserType {
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

func (u *AuthContext) EventKey() eventkey.EventKey {
	if u.UserID != uuid.Nil {
		return eventkey.NewEventKeyUserID(u.UserID)
	}
	return eventkey.NewEventKeySessionID(u.AnonSessionID)
}

func (u *AuthContext) IsGuest() bool {
	return u.UserType() == dtypes.UserTypeGuest
}

func (u *AuthContext) IsAnonymous() bool {
	return u.UserType() == dtypes.UserTypeAnonymous
}

func (u *AuthContext) IsUser() bool {
	return u.UserType() == dtypes.UserTypeUser
}

func (u *AuthContext) IsAdmin() bool {
	return u.UserType() == dtypes.UserTypeAdmin
}
