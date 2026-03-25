package authdto

import (
	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type UserContext struct {
	UserID uuid.UUID
	Roles  []dtypes.UserRole
	Token  *Token
}
