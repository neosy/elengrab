package authmw

import (
	"fmt"

	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (m *AuthMiddleware) checkHTTPSRequired(ctx *fasthttp.RequestCtx) error {
	if !nfasthttp.IsForwardedHTTPS(ctx) && m.appMode > dtypes.AppModePublic {
		return errorx.NewHTTPMessage(
			fmt.Sprintf("HTTPS is required for mode '%v'", m.appMode),
			fasthttp.StatusUpgradeRequired,
		)
	}
	return nil
}

func (m *AuthMiddleware) resolveUser(ctx *fasthttp.RequestCtx) *dauth.UserContext {
	ctxUserIface := ctx.UserValue(UserKey)
	if ctxUserIface == nil {
		return nil
	}

	userCtx, ok := ctxUserIface.(dauth.UserContext)
	if !ok {
		return nil
	}

	return &userCtx
}
