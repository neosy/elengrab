package mappers

import (
	authdto "github.com/neosy/elengrab/internal/app/usecases/auth/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

func (m *Mappers) MapUserSessionDomainToToken(
	session *dauth.UserSession,
	needsRefresh func(*dauth.UserSession) bool,
) *authdto.Token {
	return &authdto.Token{
		Token:        session.SessionToken,
		ExpiresAt:    session.ExpiresAt,
		NeedsRefresh: needsRefresh(session),
	}
}
