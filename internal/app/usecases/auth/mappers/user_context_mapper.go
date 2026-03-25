package mappers

import (
	authdto "github.com/neosy/elengrab/internal/app/usecases/auth/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

func (m *Mappers) MapUserSessionDomainToUserContext(
	user *dauth.User,
	session *dauth.UserSession,
	needsRefresh func(*dauth.UserSession) bool,
) *authdto.UserContext {
	return &authdto.UserContext{
		UserID: user.UserID,
		Roles:  user.Roles,
		Token:  m.MapUserSessionDomainToToken(session, needsRefresh),
	}
}
