package mappers

import (
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	eauth "github.com/neosy/elengrab/internal/repository/sqlite/auth/entity"
)

func (m *Mappers) MapUserSessionDomainToEntity(session *dauth.UserSession) (*eauth.UserSession, error) {
	return &eauth.UserSession{
		SessionID:    session.SessionID,
		UserID:       session.UserID,
		SessionToken: session.SessionToken,
		ExpiresAt:    session.ExpiresAt,
	}, nil
}

func (m *Mappers) MapUserSessionEntityToDomain(session *eauth.UserSession) (*dauth.UserSession, error) {
	return &dauth.UserSession{
		SessionID:    session.SessionID,
		UserID:       session.UserID,
		SessionToken: session.SessionToken,
		CreatedAt:    session.CreatedAt,
		ExpiresAt:    session.ExpiresAt,
	}, nil
}
