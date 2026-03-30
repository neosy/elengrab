package authsession

import (
	"context"

	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (uc *UserSession) FindBySessionID(ctx context.Context, sessionID uuid.UUID) (*dauth.UserSession, error) {
	if sessionID == uuid.Nil {
		return nil, nil
	}

	user, err := uc.userSessionRep.FindBySessionID(ctx, sessionID)
	if err != nil {
		uc.logger.Warn("Failed get user session", "error", err)
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	return user, nil
}

// GetBySessionID
// Record MUST exist — otherwise NOT_FOUND
func (uc *UserSession) GetBySessionID(ctx context.Context, sessionID uuid.UUID) (*dauth.UserSession, error) {
	user, err := uc.FindBySessionID(ctx, sessionID)
	if err != nil {
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	if user == nil {
		uc.logger.Debug("User session not found", "sessionID", sessionID)
		return nil, errorx.New("user session not found", exceptionx.NOT_FOUND)
	}

	return user, nil
}

// FindByToken
func (uc *UserSession) FindByToken(ctx context.Context, token string) (*dauth.UserSession, error) {
	if token == "" {
		return nil, nil
	}

	user, err := uc.userSessionRep.FindByToken(ctx, token)
	if err != nil {
		uc.logger.Warn("Failed to find user session", "error", err)
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	return user, nil
}

func (uc *UserSession) GetByToken(ctx context.Context, token string) (*dauth.UserSession, error) {
	session, err := uc.FindByToken(ctx, token)
	if err != nil {
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	if session == nil {
		uc.logger.Debug("User session not found", "token", token)
		return nil, errorx.New("user session not found", exceptionx.NOT_FOUND)
	}

	return session, nil
}
