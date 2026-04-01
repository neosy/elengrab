package authmw

import (
	"errors"
	"fmt"

	"github.com/neosy/elengrab/internal/app/usecases/auth"
	udto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/nfasthttp"
	"github.com/valyala/fasthttp"
)

// RequireAuth is a middleware that handles user session management.
// It checks for a session token in the request cookies, validates it,
// refreshes the session if needed, and registers a guest session if no valid session exists.
// After processing, it stores user information in the request context and calls the next handler.
func (a *AuthMiddleware) RequireAuth(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		userCtx, err := a.processAuth(ctx, a.appMode.IsGuestAllowed())
		if err != nil {
			err = fmt.Errorf("Internal Server Error: %w", err)
			nfasthttp.WriteErrorx(ctx, errorx.NewFromError(err, errorx.HttpStatusArg(fasthttp.StatusInternalServerError)))
			return
		}

		if userCtx == nil {
			next(ctx)
			return
		}

		ctx.SetUserValue(userKey, *userCtx)
		next(ctx)
	}
}

// OptionalAuth authenticates the user if a valid token is present.
// It updates the token if needed, but does not create a new user.
// Requests without a valid token continue as anonymous.
func (a *AuthMiddleware) OptionalAuth(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		userCtx, err := a.processAuth(ctx, false)
		if err != nil {
			err = fmt.Errorf("Internal Server Error: %w", err)
			nfasthttp.WriteErrorx(ctx, errorx.NewFromError(err, errorx.HttpStatusArg(fasthttp.StatusInternalServerError)))
			return
		}

		if userCtx == nil {
			next(ctx)
			return
		}

		ctx.SetUserValue(userKey, *userCtx)
		next(ctx)
	}
}

func (a *AuthMiddleware) processAuth(ctx *fasthttp.RequestCtx, createGuest bool) (*dauth.UserContext, error) {
	var session *udto.AuthUserResponse

	token := CookieSessionTokenKey.GetValue(ctx)
	if token != "" {
		var err error
		session, err = a.auth.ValidateSession(ctx, token)
		if err != nil {
			switch {
			case errors.Is(err, auth.ErrSessionNotFound):
			case errors.Is(err, auth.ErrSessionExpired):
			case errors.Is(err, auth.ErrUserNotFound):
			default:
				return nil, err
			}
		}

		if session != nil && session.Token.NeedsRefresh {
			newSession, err := a.auth.RefreshSession(ctx, session.Token.Token)
			if err != nil {
				a.logger.Warn("failed to refresh session: %v", "error", err)
				return nil, err
			}
			if newSession.Token.Token != session.Token.Token {
				session = newSession
				CookieSessionTokenKey.SetValue(ctx, session.Token.Token, session.Token.ExpiresAt)
			}
		}
	}

	if session == nil && createGuest {
		var err error
		session, err = a.auth.RegisterGuest(ctx)
		if err != nil {
			return nil, err
		}
		CookieSessionTokenKey.SetValue(ctx, session.Token.Token, session.Token.ExpiresAt)
	}

	var userCtx *dauth.UserContext

	if session != nil {
		userCtx = &dauth.UserContext{
			UserID: session.UserID,
			Login:  session.Login,
			Roles:  session.Roles,
		}
	}

	return userCtx, nil
}
