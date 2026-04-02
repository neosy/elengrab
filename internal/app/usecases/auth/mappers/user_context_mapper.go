package mappers

import (
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

func (m *Mappers) MapUserSessionDomainToUserContext(
	user *dauth.User,
	session *dauth.UserSession,
	needsRefresh func(*dauth.UserSession) bool,
) *dto.AuthUserResponse {
	return &dto.AuthUserResponse{
		UserID: user.UserID,
		Login:  user.Login.String(),
		Email:  uptr.Deref(user.Email),
		Roles:  user.Roles,
		Token:  m.MapUserSessionDomainToToken(session, needsRefresh),
	}
}
