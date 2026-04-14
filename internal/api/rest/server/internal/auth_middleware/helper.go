package authmw

import (
	"fmt"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
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
