package authmiddleware

import (
	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/valyala/fasthttp"
)

// Authorize middleware checks the session token and sets user_id in context if a valid session exists.
// Does NOT create a new user or session.
func (a *AuthMiddleware) Authorize(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		token := a.getSessionToken(ctx)
		if token == "" {
			next(ctx)
			return
		}

		session, err := a.auth.GetSessionByToken(ctx, token)
		if err != nil {
			next(ctx)
			return
		}

		ctx.SetUserValue(userIDKey, session.UserID)
		next(ctx)
	}
}

// AutoRegister middleware checks the session token and authorizes the user.
// If no user exists for the token, it automatically creates a new user and a session,
// then sets user_id in the context.
func (a *AuthMiddleware) AutoRegister(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var (
			session *dauth.UserSession
			userID  uuid.UUID
			err     error
		)

		token := a.getSessionToken(ctx)
		if token != "" {
			session, err = a.auth.GetSessionByToken(ctx, token)
			if err != nil {
				cookieSessionTokenKey.deleteCookie(ctx)
				if session == nil {
					next(ctx)
					return
				}
				userID = session.UserID
				session = nil
			} else {
				userID = session.UserID
			}
		}

		if userID == uuid.Nil {
			user, err := a.auth.CreateUser(ctx)
			if err != nil {
				next(ctx)
				return
			}
			userID = user.UserID
		}

		if session == nil {
			session, err = a.createSession(ctx, userID)
			if err != nil {
				next(ctx)
				return
			}
		}

		ctx.SetUserValue(userIDKey, session.UserID)
		next(ctx)
	}
}
