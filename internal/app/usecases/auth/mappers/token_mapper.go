package mappers

import (
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (m *Mappers) MapUserSessionDomainToToken(
	session *dauth.UserSession,
	needsRefresh func(*dauth.UserSession) bool,
) *dtypes.AuthToken {
	var nr bool
	if needsRefresh != nil {
		nr = needsRefresh(session)
	}
	return &dtypes.AuthToken{
		Token:        session.SessionToken,
		ExpiresAt:    session.ExpiresAt,
		NeedsRefresh: nr,
	}
}
