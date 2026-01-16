package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/pkg/errorx"
	"github.com/neosy/elengrab/pkg/errorx/exceptionx"
)

func (u *Auth) CreateSession(ctx context.Context, userID uuid.UUID) (*dauth.UserSession, error) {
	token, err := u.generateSessionToken()
	if err != nil {
		return nil, errorx.NewByErr(err, exceptionx.ERROR)
	}

	session := &dauth.UserSession{
		UserID:       userID,
		SessionToken: token,
		ExpiresAt:    time.Now().Add(sessionTTL),
	}

	sessionID, err := u.userSession.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	session, err = u.userSession.GetBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (u *Auth) generateSessionToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
