package mappers

import (
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

func (m *Mappers) MapUserToUserResponse(user *dauth.User) *dto.AuthUserResponse {
	return &dto.AuthUserResponse{
		UserID:  user.UserID,
		Login:   user.Login.String(),
		Email:   uptr.Deref(user.Email),
		RoleIDs: user.RoleIDs,
	}
}
