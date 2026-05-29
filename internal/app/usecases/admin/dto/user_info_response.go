package dto

import (
	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

type UserInfoResponse struct {
	User                dauth.User
	RolesWithAssignment []RoleWithAssignment
}
