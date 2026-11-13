package pservices

import (
	"uuid"

	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	ierrors "github.com/neosy/elengrab/internal/errors"
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

	Token *dtypes.AuthToken
}

func (u *AuthUserResponse) UserType() dtypes.UserType {
	roleIDs, _ := dtypes.ParseUserRoleIDs(u.RoleIDs)
	return dauth.ResolveUserType(u.UserID, roleIDs)
}

// User is a user returned by the service
func (r *AuthUserResponse) User() *dauth.User {
	if r == nil {
		return nil
	}

	var email *string
	if r.Email != "" {
		email = &r.Email
	}

	return &dauth.User{
		UserID:  r.UserID,
		Login:   dtypes.Login(r.Login),
		Email:   email,
		RoleIDs: r.RoleIDs,
	}
}

// AuthContext converts the authentication user response into an authentication context.
func (r *AuthUserResponse) AuthContext(guestCreated bool) (*dauth.AuthContext, error) {
	if r == nil {
		return nil, ierrors.ErrMethodReceiverNullPointer
	}

	roleIDs, err := dtypes.ParseUserRoleIDs(r.RoleIDs)
	if err != nil {
		return nil, err
	}

	return &dauth.AuthContext{
		UserID:       r.UserID,
		Login:        r.Login,
		Email:        r.Email,
		RoleIDs:      roleIDs,
		GuestCreated: guestCreated,
	}, nil
}

// HasDifferentToken reports whether the authentication token differs from another response.
func (r *AuthUserResponse) HasDifferentToken(other *AuthUserResponse) bool {
	return r != nil && other != nil &&
		r.Token != nil && other.Token != nil &&
		r.Token.Token != other.Token.Token
}
