package persistence

import (
	"context"

	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

type UserSessionRepositoryFactory func() UserSessionRepository

type UserSessionRepository interface {
	Transactional

	Insert(ctx context.Context, session *dauth.UserSession) error
	Update(ctx context.Context, session *dauth.UserSession) error
	Save(ctx context.Context, session *dauth.UserSession) error
	FindBySessionID(ctx context.Context, sessionID uuid.UUID) (*dauth.UserSession, error)
	FindByToken(ctx context.Context, token string) (*dauth.UserSession, error)
}
