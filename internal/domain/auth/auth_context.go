package dauth

import (
	"uuid"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	eventkey "github.com/neosy/elengrab/internal/domain/types/event_key"
)

type AuthContext struct {
	UserID        uuid.UUID
	AnonSessionID uuid.UUID

	Login string
	Email string

	RoleIDs dtypes.UserRoleIDs

	GuestCreated bool
}

// NewEmptyAuthContext creates an empty authentication context.
func NewEmptyAuthContext() AuthContext {
	return AuthContext{}
}

// NewAnonymousAuthContext creates a new AuthContext for an anonymous user.
func NewAnonymousAuthContext() AuthContext {
	return AuthContext{
		UserID:  AnonymousUserID(),
		RoleIDs: nil,
	}
}

// NewAnonymousSessionAuthContext creates a new AuthContext for an anonymous session.
func NewAnonymousSessionAuthContext(anonSessionID uuid.UUID) AuthContext {
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

// EventKey returns the event key for the user.
func (u *AuthContext) EventKey() eventkey.EventKey {
	if IsRegisteredUserID(u.UserID) {
		return eventkey.NewEventKeyUserID(u.UserID)
	}

	return eventkey.NewEventKeySessionID(u.AnonSessionID)
}

// IsAnonymous checks whether the authentication context belongs to an anonymous user.
func (u *AuthContext) IsAnonymous() bool {
	return u.UserType() == dtypes.UserTypeAnonymous
}

// IsAnonymousWithoutSession checks whether the anonymous user has no session.
func (u *AuthContext) IsAnonymousWithoutSession() bool {
	return u.IsAnonymous() && u.AnonSessionID == uuid.Nil()
}

// IsAnonymousWithSession checks whether the anonymous user has a session.
func (u *AuthContext) IsAnonymousWithSession() bool {
	return u.IsAnonymous() && u.AnonSessionID != uuid.Nil()
}

// IsGuest reports whether the authentication context belongs to a guest user.
func (u *AuthContext) IsGuest() bool {
	return u.UserType() == dtypes.UserTypeGuest
}

// IsUser reports whether the authentication context belongs to a user.
func (u *AuthContext) IsUser() bool {
	return u.UserType() == dtypes.UserTypeUser
}

// IsAuthenticatedUser reports whether the authentication context belongs to an authenticated user.
func (u *AuthContext) IsAuthenticatedUser() bool {
	return u.UserType() >= dtypes.UserTypeUser
}

// IsGuestOrHigher reports whether the authentication context belongs to a guest or higher user type.
func (u *AuthContext) IsGuestOrHigher() bool {
	return u.UserType() >= dtypes.UserTypeGuest
}

// IsAdmin reports whether the authentication context belongs to an admin user.
func (u *AuthContext) IsAdmin() bool {
	return u.UserType() == dtypes.UserTypeAdmin
}

// IsGuestOrUser reports whether the authentication context belongs to a guest or regular user.
func (u *AuthContext) IsGuestOrUser() bool {
	return u.IsUser() || u.IsGuest()
}

// HasSession checks whether the authentication context has an identifiable session.
func (u *AuthContext) HasSession() bool {
	return u.IsGuestOrHigher() || u.IsAnonymousWithSession()
}
