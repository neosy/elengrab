package dto

import (
	"github.com/google/uuid"
)

type AuthUserResponse struct {
	UserID  uuid.UUID
	Login   string
	Email   string
	RoleIDs []string
	Token   *AuthToken
}
