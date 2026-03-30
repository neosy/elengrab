package authmw

import (
	"fmt"
	"time"

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

func (k cookieKey) SetValue(ctx *fasthttp.RequestCtx, value any, expiresAt time.Time) {
	ctx.Response.Header.SetCookie(k.makeCookie(value, &expiresAt, nil))
}

func (k cookieKey) Delete(ctx *fasthttp.RequestCtx) {
	past := time.Now().Add(-time.Hour)
	ctx.Response.Header.SetCookie(k.makeCookie("", &past, nil))
}

func (k cookieKey) compareValue(ctx *fasthttp.RequestCtx, value any) bool {
	v := string(ctx.Request.Header.Cookie(string(k)))
	return v == fmt.Sprint(value)
}

func (k cookieKey) GetValue(ctx *fasthttp.RequestCtx) string {
	return string(ctx.Request.Header.Cookie(string(k)))
}
