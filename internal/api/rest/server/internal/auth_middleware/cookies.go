package authmiddleware

import (
	"fmt"
	"time"

	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
	"github.com/valyala/fasthttp"
)

type cookieKey string

func (k cookieKey) makeCookie(value any, expiresAt *time.Time, maxAge *int) *fasthttp.Cookie {
	var cookie fasthttp.Cookie

	cookie.SetKey(string(k))
	cookie.SetValue(fmt.Sprint(value))
	cookie.SetPath("/")
	cookie.SetSameSite(fasthttp.CookieSameSiteLaxMode)
	if expiresAt != nil {
		cookie.SetExpire(*expiresAt)
	}
	if maxAge != nil {
		cookie.SetMaxAge(*maxAge)
	}

	return &cookie
}

func (k cookieKey) setCookie(ctx *fasthttp.RequestCtx, value any, expiresAt time.Time) {
	ctx.Response.Header.SetCookie(k.makeCookie(value, &expiresAt, nil))
}

func (k cookieKey) deleteCookie(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetCookie(k.makeCookie("", nil, uptr.Any(-1)))
}

func (k cookieKey) compareValue(ctx *fasthttp.RequestCtx, value any) bool {
	v := string(ctx.Request.Header.Cookie(string(k)))
	return v == fmt.Sprint(value)
}

func (k cookieKey) getValue(ctx *fasthttp.RequestCtx) string {
	return string(ctx.Request.Header.Cookie(string(k)))
}
