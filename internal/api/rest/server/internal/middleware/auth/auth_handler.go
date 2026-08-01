package authmw

import (
	"errors"

	"github.com/google/uuid"
	autherr "github.com/neosy/elengrab/internal/app/usecases/auth/errors"
	udto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

// AuthOrGuest is a middleware that manages user sessions.
// It checks for a session token in the request cookies and validates it.
// If the token is valid, it refreshes the session as needed.
// If no valid session exists, it may create a guest session depending on the application mode.
// Finally, it stores the user information (authenticated or guest) in the request context and calls the next handler.
func (a *AuthMiddleware) AuthOrGuest(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if err := a.checkHTTPS(ctx); err != nil {
			nfasthttp.WriteErrorx(ctx, err)
			return
		}

		userCtx, err := a.processAuth(ctx, a.appMode.IsGuestAllowed())
		if err != nil {
			nfasthttp.WriteErrorx(
				ctx,
				errorx.Errorf(
					"internal Server Error: %w", err,
					errorx.WithHttpStatus(fasthttp.StatusInternalServerError),
				),
			)
			return
		}

		// If no valid session exists and guest sessions are not allowed,
		// return an unauthorized response.
		if userCtx == nil && !a.appMode.IsUserRequiredForWrite() {
			a.SetCookieAnonym(ctx)
			next(ctx)
			return
		}

		if userCtx == nil {
			nfasthttp.WriteErrorx(ctx, ierrors.ErrUnauthorized)
			return
		}

		ctx.SetUserValue(UserKey, *userCtx)
		next(ctx)
	}
}

// RequireAuth is a middleware that enforces strict user authentication.
// It checks for a session token in the request cookies and validates it.
// If the token is valid, it refreshes the session as needed.
// If no valid session exists, it returns an unauthorized (401) response.
// Finally, it stores the authenticated user information in the request context
// and calls the next handler.
func (a *AuthMiddleware) RequireAuth(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if err := a.checkHTTPS(ctx); err != nil {
			nfasthttp.WriteErrorx(ctx, err)
			return
		}

		userCtx, err := a.processAuth(ctx, false)
		if err != nil {
			nfasthttp.WriteErrorx(
				ctx,
				errorx.Errorf(
					"internal Server Error: %w", err,
					errorx.WithHttpStatus(fasthttp.StatusInternalServerError),
				),
			)
			return
		}

		if userCtx == nil {
			nfasthttp.WriteErrorx(ctx, ierrors.ErrUnauthorized)
			return
		}

		ctx.SetUserValue(UserKey, *userCtx)
		next(ctx)
	}
}

// RequireAdmin is a middleware that allows access only to authenticated
// users with administrator privileges. It returns 401 if the user is not
// authenticated and 403 if the user is not an administrator.
func (a *AuthMiddleware) RequireAdmin(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return a.RequireAuth(func(ctx *fasthttp.RequestCtx) {
		user := a.resolveUser(ctx)
		if user == nil {
			nfasthttp.WriteErrorx(
				ctx,
				ierrors.ErrUnauthorized,
			)
			return
		}

		if user.UserType() != dtypes.UserTypeAdmin {
			nfasthttp.WriteErrorx(
				ctx,
				errorx.Errorf(
					"admin privileges required",
					errorx.WithHttpStatus(fasthttp.StatusForbidden),
				),
			)
			return
		}

		next(ctx)
	})
}

// AuthOrAnonym authenticates the user if a valid token is present.
// It updates the token if needed, but does not create a new user.
// Requests without a valid token continue as anonymous.
func (a *AuthMiddleware) AuthOrAnonym(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		userCtx, err := a.processAuth(ctx, false)
		if err != nil {
			nfasthttp.WriteErrorx(
				ctx,
				errorx.Errorf(
					"internal Server Error: %w", err,
					errorx.WithHttpStatus(fasthttp.StatusInternalServerError),
				),
			)
			return
		}

		if userCtx == nil {
			a.SetCookieAnonym(ctx)
			next(ctx)
			return
		}

		ctx.SetUserValue(UserKey, *userCtx)
		next(ctx)
	}
}

// AuthOptional is a middleware that optionally enforces user authentication.
// It checks for a session token in the request cookies and validates it.
// If the token is valid, it refreshes the session as needed.
// If no valid session exists, its behavior depends on the application mode:
//   - In strict authentication mode (AppModeAuthOnly), it returns an unauthorized (401) response.
//   - Otherwise, it treats the user as anonymous (nil) and proceeds to the next handler.
//
// Finally, if a user (authenticated) is present, it stores the user information
// in the request context before calling the next handler.
func (a *AuthMiddleware) AuthOptional(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		userCtx, err := a.processAuth(ctx, false)
		if err != nil {
			nfasthttp.WriteErrorx(
				ctx,
				errorx.Errorf(
					"internal Server Error: %w", err,
					errorx.WithHttpStatus(fasthttp.StatusInternalServerError),
				),
			)
			return
		}

		if userCtx == nil {
			if a.appMode == dtypes.AppModeAuthenticated {
				nfasthttp.WriteErrorx(ctx, ierrors.ErrUnauthorized)
				return
			}

			a.SetCookieAnonym(ctx)
			next(ctx)
			return
		}

		ctx.SetUserValue(UserKey, *userCtx)
		next(ctx)
	}
}

// RequireAuthMode is a middleware that checks user authentication based on the application mode.
// - If the user is authenticated, their information is stored in the request context.
// - If no valid session exists and the current app mode requires a user, it returns 401 Unauthorized.
// - If no valid session exists but the app mode does not require a user, it calls the next handler with nil (anonymous) user.
func (a *AuthMiddleware) RequireAuthMode(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if err := a.checkHTTPS(ctx); err != nil {
			nfasthttp.WriteErrorx(ctx, err)
			return
		}

		userCtx, err := a.processAuth(ctx, false)
		if err != nil {
			nfasthttp.WriteErrorx(
				ctx,
				errorx.Errorf(
					"internal Server Error: %w", err,
					errorx.WithHttpStatus(fasthttp.StatusInternalServerError),
				),
			)
			return
		}

		if userCtx == nil {
			if a.appMode.IsUserRequiredForWrite() {
				nfasthttp.WriteErrorx(ctx, ierrors.ErrUnauthorized)
				return
			}

			a.SetCookieAnonym(ctx)
			next(ctx)
			return
		}

		ctx.SetUserValue(UserKey, *userCtx)
		next(ctx)
	}
}

func (a *AuthMiddleware) processAuth(ctx *fasthttp.RequestCtx, createGuest bool) (*dauth.AuthContext, error) {
	var session *udto.AuthUserResponse

	if !nfasthttp.IsForwardedHTTPS(ctx) {
		a.SetCookieAnonym(ctx)
		return nil, nil
	}

	token := CookieSessionTokenKey.GetValue(ctx)

	if token != "" {
		var err error
		session, err = a.auth.ValidateSession(ctx, token)
		if err != nil {
			switch {
			case errors.Is(err, autherr.ErrSessionNotFound):
			case errors.Is(err, autherr.ErrSessionExpired):
			case errors.Is(err, autherr.ErrUserNotFound):
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
				a.SetCookieToken(ctx, session.Token.Token, WithExpiresAt(session.Token.ExpiresAt))
			}
		}
	}

	var guestCreated bool
	if session == nil && createGuest {
		var err error
		session, err = a.auth.RegisterGuest(ctx)
		if err != nil {
			return nil, err
		}
		guestCreated = true
		a.SetCookieToken(ctx, session.Token.Token, WithExpiresAt(session.Token.ExpiresAt))
	}

	if session != nil {
		roleIDs, err := dtypes.ParseUserRoleIDs(session.RoleIDs)
		if err != nil {
			return nil, err
		}
		userCtx := &dauth.AuthContext{
			UserID:       session.UserID,
			Login:        session.Login,
			Email:        session.Email,
			RoleIDs:      roleIDs,
			GuestCreated: guestCreated,
		}
		return userCtx, nil
	}

	return nil, nil
}

func (m *AuthMiddleware) SetCookieAnonym(ctx *fasthttp.RequestCtx, opts ...cookieOption) {
	if CookieAnonSessionIDKey.GetValue(ctx) != "" {
		return
	}
	CookieAnonSessionIDKey.SetValue(ctx, uuid.New().String(), append([]cookieOption{WithHTTPOnly()}, opts...)...)
}

func (m *AuthMiddleware) SetCookieToken(ctx *fasthttp.RequestCtx, token string, opts ...cookieOption) {
	CookieSessionTokenKey.SetValueWithSecure(ctx, token, opts...)
}
