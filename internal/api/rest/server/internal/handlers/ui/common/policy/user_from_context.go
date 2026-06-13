package policy

import (
	"github.com/google/uuid"
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	authmw "github.com/neosy/elengrab/internal/api/rest/server/internal/middleware/auth"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

// ResolveUser extracts authenticated user from fasthttp request context.
// It returns nil if no user was set by authentication middleware,
// if the value has an invalid type, or if the user ID is empty.
func ResolveUser(ctx *fasthttp.RequestCtx) *dauth.UserContext {
	ctxUserIface := ctx.UserValue(authmw.UserKey)
	if ctxUserIface == nil {
		return nil
	}

	userCtx, ok := ctxUserIface.(dauth.UserContext)
	if !ok {
		return nil
	}

	if userCtx.UserID == dauth.AnonymousUserID() {
		return nil
	}

	return &userCtx
}

// ResolveUserOrAnonym returns authenticated user from context if present,
// otherwise returns anonymous user context.
func ResolveUserOrAnonym(ctx *fasthttp.RequestCtx) *dauth.UserContext {
	userCtx := ResolveUser(ctx)

	if userCtx != nil {
		return userCtx
	}

	var anonSessionID uuid.UUID
	if id := authmw.CookieAnonSessionIDKey.GetValue(ctx); id != "" {
		anonSessionID, _ = uuid.Parse(id)
	}

	return dauth.UserContextAnonymous(anonSessionID)
}

// EnsureUser returns authenticated user from fasthttp request context.
// If no user is present (or user is invalid), it returns HTTP unauthorized error.
func EnsureUser(ctx *fasthttp.RequestCtx) (*dauth.UserContext, error) {
	userCtx := ResolveUser(ctx)

	if userCtx == nil && !nfasthttp.IsForwardedHTTPS(ctx) {
		return nil, errorx.NewHTTPMessage("HTTPS is required for authentication", fasthttp.StatusBadRequest)
	}

	if userCtx == nil {
		return nil, apierrors.ErrUnauthorized
	}

	return userCtx, nil
}

// ResolveUserOrFallback resolves user from request context applying access policy:
// in public app mode over non-HTTPS requests it returns anonymous user,
// otherwise it requires authenticated user and returns error if absent.
func ResolveUserOrFallback(ctx *fasthttp.RequestCtx, appMode dtypes.AppMode) (*dauth.UserContext, error) {
	if IsAnonymousFallbackAllowed(ctx, appMode) {
		return ResolveUserOrAnonym(ctx), nil
	}
	return EnsureUser(ctx)
}

// IsAnonymousFallbackAllowed returns true if anonymous fallback is allowed for current request.
// It's used in public app mode over non-HTTPS requests.
func IsAnonymousFallbackAllowed(ctx *fasthttp.RequestCtx, appMode dtypes.AppMode) bool {
	return !nfasthttp.IsForwardedHTTPS(ctx) && !appMode.IsUserRequired()
}
