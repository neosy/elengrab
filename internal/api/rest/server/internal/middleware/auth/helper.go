package authmw

import (
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (m *AuthMiddleware) validateHTTPSPolicy() bool {
	return m.appMode > dtypes.AppModePublic
}

func (m *AuthMiddleware) requireHTTPS(ctx *fasthttp.RequestCtx) error {
	if !nfasthttp.IsForwardedHTTPS(ctx) {
		return apierrors.ErrHTTPSRequired
	}
	return nil
}

func (m *AuthMiddleware) checkHTTPS(ctx *fasthttp.RequestCtx) error {
	if m.validateHTTPSPolicy() {
		return m.requireHTTPS(ctx)
	}
	return nil
}

func (m *AuthMiddleware) resolveUser(ctx *fasthttp.RequestCtx) *dauth.AuthContext {
	ctxUserIface := ctx.UserValue(UserKey)
	if ctxUserIface == nil {
		return nil
	}

	userCtx, ok := ctxUserIface.(dauth.AuthContext)
	if !ok {
		return nil
	}

	return &userCtx
}
