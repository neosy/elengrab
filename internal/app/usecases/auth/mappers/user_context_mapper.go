package mappers

import (
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

func (m *Mappers) MapUserSessionDomainToUserContext(
	user *dauth.User,
	session *dauth.UserSession,
	needsRefresh func(*dauth.UserSession) bool,
) *dto.AuthUserResponse {
	return &dto.AuthUserResponse{
		UserID: user.UserID,
		Roles:  user.Roles,
		Token:  m.MapUserSessionDomainToToken(session, needsRefresh),
	}
}
