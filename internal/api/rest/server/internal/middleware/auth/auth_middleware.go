package authmw

import (
	"log/slog"
	"net/netip"

	iconfig "github.com/neosy/elengrab/internal/config"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

type AuthMiddleware struct {
	logger *slog.Logger

	// services
	auth pservices.AuthMiddleware

	// options
	appMode      dtypes.AppMode
	reverseProxy authMiddlewareReverseProxy
}

type authMiddlewareReverseProxy struct {
	enableAuth     bool
	trustedProxies []netip.Prefix
}

func (r authMiddlewareReverseProxy) isReverseProxyAuthEnabled() bool {
	return r.enableAuth && len(r.trustedProxies) > 0
}

func NewAuthMiddleware(
	logger *slog.Logger,
	auth pservices.AuthMiddleware,
	appMode dtypes.AppMode,
	reverseProxy iconfig.HTTPServerReverseProxyConfig,
) *AuthMiddleware {
	trustedProxies, err := parseTrustedProxies(reverseProxy.TrustedProxies)
	if err != nil {
		logger.Warn("Failed to parse trusted reverse proxy addresses", "error", err)
	}

	return &AuthMiddleware{
		logger: logger,

		// services
		auth: auth,

		// options
		appMode: appMode,
		reverseProxy: authMiddlewareReverseProxy{
			enableAuth:     reverseProxy.EnableAuth,
			trustedProxies: trustedProxies,
		},
	}
}
