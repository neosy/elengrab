package dto

import (
	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type AuthUserResponse struct {
	UserID uuid.UUID
	Login  string
	Roles  []dtypes.UserRole
	Token  *AuthToken
}
