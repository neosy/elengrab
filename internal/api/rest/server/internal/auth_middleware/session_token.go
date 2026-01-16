package authmiddleware

import (
	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/valyala/fasthttp"
)

func (a *AuthMiddleware) getSessionToken(ctx *fasthttp.RequestCtx) string {
	return cookieSessionTokenKey.getValue(ctx)
}

func (a *AuthMiddleware) createSession(ctx *fasthttp.RequestCtx, userID uuid.UUID) (*dauth.UserSession, error) {
	session, err := a.auth.CreateSession(ctx, userID)
	if err != nil {
		return nil, err
	}

	cookieSessionTokenKey.setCookie(ctx, session.SessionToken, session.ExpiresAt)

	return session, nil
}
