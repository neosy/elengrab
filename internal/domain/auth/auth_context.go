package dauth

import (
	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	eventkey "github.com/neosy/elengrab/internal/domain/types/event_key"
)

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
	if u == nil {
		return dtypes.UserTypeAnonymous
	}

	return ResolveUserType(u.UserID, u.RoleIDs)
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

// IsRegularUser reports whether the authentication context belongs to a regular or guest user.
// Anonymous and admin authentication contexts are not considered regular users.
func (u *AuthContext) IsRegularUser() bool {
	return u.IsUser() || u.IsGuest()
}
