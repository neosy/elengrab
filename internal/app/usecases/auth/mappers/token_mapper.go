package mappers

import (
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

func (m *Mappers) MapUserSessionDomainToToken(
	session *dauth.UserSession,
	needsRefresh func(*dauth.UserSession) bool,
) *dto.AuthToken {
	var nr bool
	if needsRefresh != nil {
		nr = needsRefresh(session)
	}
	return &dto.AuthToken{
		Token:        session.SessionToken,
		ExpiresAt:    session.ExpiresAt,
		NeedsRefresh: nr,
	}
}
