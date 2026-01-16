package mappers

import (
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
)

func (m *Mappers) MapUserSessionDomainToEntity(session *dauth.UserSession) (*edownload.UserSession, error) {
	return &edownload.UserSession{
		SessionID:    session.SessionID,
		UserID:       session.UserID,
		SessionToken: session.SessionToken,
		ExpiresAt:    session.ExpiresAt,
	}, nil
}

func (m *Mappers) MapUserSessionEntityToDomain(session *edownload.UserSession) (*dauth.UserSession, error) {
	return &dauth.UserSession{
		SessionID:    session.SessionID,
		UserID:       session.UserID,
		SessionToken: session.SessionToken,
		CreatedAt:    session.CreatedAt,
		ExpiresAt:    session.ExpiresAt,
	}, nil
}
