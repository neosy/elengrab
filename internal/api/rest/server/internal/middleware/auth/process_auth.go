package authmw

import (
	"errors"
	"fmt"
	"net/netip"

	"github.com/neosy/elengrab/internal/app/usecases/auth"
	autherr "github.com/neosy/elengrab/internal/app/usecases/auth/errors"
	authweb "github.com/neosy/elengrab/internal/app/usecases/auth_web"
	udto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	pservices "github.com/neosy/elengrab/internal/ports/services"
	"github.com/valyala/fasthttp"
)

func (a *AuthMiddleware) processAuth(ctx *fasthttp.RequestCtx, createGuest bool) (*dauth.AuthContext, error) {
	authCtx, err := a.processReverseProxyAuth(ctx)
	if err != nil {
		return nil, err
	}
	if authCtx != nil {
		return authCtx, nil
	}

	var authUser *udto.AuthUserResponse

	if !nfasthttp.IsForwardedHTTPS(ctx) {
		a.SetCookieAnonym(ctx)
		return nil, nil
	}

	token := CookieSessionTokenKey.GetValue(ctx)

	if token != "" {
		var err error
		authUser, err = a.auth.ValidateSession(ctx, token)
		if err != nil {
			switch {
			case errors.Is(err, autherr.ErrSessionNotFound):
			case errors.Is(err, autherr.ErrSessionExpired):
			case errors.Is(err, autherr.ErrUserNotFound):
			default:
				return nil, err
			}
		}

		if authUser != nil && authUser.Token.NeedsRefresh {
			refreshedAuthUser, err := a.auth.RefreshSession(ctx, authUser.Token.Token)
			if err != nil {
				a.logger.Warn("Failed to refresh session", "error", err)
				return nil, err
			}
			if refreshedAuthUser.HasDifferentToken(authUser) {
				authUser = refreshedAuthUser
				a.SetCookieToken(ctx, authUser.Token.Token, WithExpiresAt(authUser.Token.ExpiresAt))
			}
		}
	}

	var guestCreated bool
	if authUser == nil && createGuest {
		authUser, err = a.auth.RegisterGuest(ctx)
		if err != nil {
			return nil, err
		}

		guestCreated = true
		a.SetCookieToken(ctx, authUser.Token.Token, WithExpiresAt(authUser.Token.ExpiresAt))
	}

	if authUser != nil {
		return authUser.AuthContext(guestCreated)
	}

	return nil, nil
}

// processReverseProxyAuth checks for a user login in the reverse proxy authentication header.
func (a *AuthMiddleware) processReverseProxyAuth(ctx *fasthttp.RequestCtx) (*dauth.AuthContext, error) {
	if !a.reverseProxy.isReverseProxyAuthEnabled() {
		return nil, nil
	}

	if !a.isTrustedProxy(ctx) {
		return nil, nil
	}

	login := string(ctx.Request.Header.Peek(HeaderAuthRequestUser))
	if login == "" {
		return nil, nil
	}

	user, err := a.auth.FindByLogin(ctx, login)
	if err != nil {
		a.logger.Warn("Failed to find user by login from reverse proxy authentication header", "login", login, "error", err)
		return nil, err
	}

	if user != nil && !user.IsActive {
		a.logger.Warn("User found by reverse proxy authentication header is inactive", "login", login)
		return nil, nil
	}

	if user == nil {
		user, err = a.registerReverseProxyUser(ctx, login)
		if err != nil || user == nil {
			return nil, err
		}
	}

	var email string
	if user.Email != nil {
		email = *user.Email
	}

	roleIDs, err := user.ParseRoleIDs()
	if err != nil {
		a.logger.Warn(
			"Failed to parse user role IDs",
			"login", login, "role_ids",
			user.RoleIDs, "error", err,
		)
		return nil, nil
	}

	return &dauth.AuthContext{
		UserID:  user.UserID,
		Login:   user.Login.String(),
		Email:   email,
		RoleIDs: roleIDs,
	}, nil
}

// registerReverseProxyUser registers a new user authenticated by a reverse proxy.
func (a *AuthMiddleware) registerReverseProxyUser(
	ctx *fasthttp.RequestCtx,
	login string,
) (*dauth.User, error) {
	if err := authweb.ValidateLogin(login); err != nil {
		a.logger.Warn("Invalid login from reverse proxy authentication header", "login", login, "error", err)
		return nil, nil
	}

	password, err := auth.GenerateRandomPassword()
	if err != nil {
		a.logger.Warn("Failed to generate random password for reverse proxy user", "login", login, "error", err)
		return nil, nil
	}

	userRegistered, err := a.auth.RegisterUser(ctx,
		&pservices.AuthRegisterUserRequest{
			Login:    login,
			Password: password,
		},
	)
	if err != nil {
		a.logger.Warn("Failed to register user from reverse proxy authentication header", "login", login, "error", err)
		return nil, err
	}

	if userRegistered == nil {
		a.logger.Warn("User registration returned a nil result", "login", login)
		return nil, nil
	}

	return userRegistered.User(), nil
}

// parseTrustedProxies parses trusted proxy IP addresses and CIDR ranges.
func parseTrustedProxies(addresses []string) ([]netip.Prefix, error) {
	prefixes := make([]netip.Prefix, 0, len(addresses))

	for _, address := range addresses {
		prefix, err := netip.ParsePrefix(address)
		if err != nil {
			ip, ipErr := netip.ParseAddr(address)
			if ipErr != nil {
				return nil, fmt.Errorf("invalid trusted proxy address %q: %w", address, ipErr)
			}

			prefix = netip.PrefixFrom(ip, ip.BitLen())
		}

		prefixes = append(prefixes, prefix.Masked())
	}

	return prefixes, nil
}

// isTrustedProxy checks whether the request comes from a trusted reverse proxy.
func (a *AuthMiddleware) isTrustedProxy(ctx *fasthttp.RequestCtx) bool {
	ip, ok := netip.AddrFromSlice(ctx.RemoteIP())
	if !ok {
		return false
	}

	ip = ip.Unmap()

	for _, prefix := range a.reverseProxy.trustedProxies {
		if prefix.Contains(ip) {
			return true
		}
	}

	return false
}
