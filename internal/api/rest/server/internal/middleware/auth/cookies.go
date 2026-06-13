package authmw

import (
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

type (
	cookieKey     string
	cookieOptions struct {
		expiresAt *time.Time
		maxAge    *int
		httpOnly  bool
		secure    bool
	}
	cookieOption func(*cookieOptions)
)

func WithExpiresAt(t time.Time) cookieOption {
	return func(o *cookieOptions) {
		o.expiresAt = &t
	}
}

func WithMaxAge(seconds int) cookieOption {
	return func(o *cookieOptions) {
		o.maxAge = &seconds
	}
}

func WithHTTPOnly() cookieOption {
	return func(o *cookieOptions) {
		o.httpOnly = true
	}
}

func WithSecure() cookieOption {
	return func(o *cookieOptions) {
		o.secure = true
	}
}

func (k cookieKey) makeCookie(value string, opts ...cookieOption) *fasthttp.Cookie {
	var cookie fasthttp.Cookie

	options := cookieOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	cookie.SetKey(string(k))
	cookie.SetValue(value)
	cookie.SetPath("/")
	cookie.SetSameSite(fasthttp.CookieSameSiteLaxMode)
	cookie.SetHTTPOnly(options.httpOnly)
	cookie.SetSecure(options.secure)
	if options.expiresAt != nil {
		cookie.SetExpire(*options.expiresAt)
	}
	if options.maxAge != nil {
		cookie.SetMaxAge(*options.maxAge)
	}

	return &cookie
}

func (k cookieKey) String() string {
	return string(k)
}

func (k cookieKey) SetValue(ctx *fasthttp.RequestCtx, value string, opts ...cookieOption) {
	ctx.Response.Header.SetCookie(k.makeCookie(value, opts...))
}

func (k cookieKey) SetValueWithSecure(ctx *fasthttp.RequestCtx, value string, opts ...cookieOption) {
	newOpts := append([]cookieOption{WithHTTPOnly(), WithSecure()}, opts...)
	k.SetValue(ctx, value, newOpts...)
}

func (k cookieKey) Delete(ctx *fasthttp.RequestCtx, opts ...cookieOption) {
	past := time.Now().Add(-time.Hour)
	newOpts := append([]cookieOption{WithExpiresAt(past), WithMaxAge(0)}, opts...)
	ctx.Response.Header.SetCookie(k.makeCookie("", newOpts...))
}

func (k cookieKey) DeleteWithSecure(ctx *fasthttp.RequestCtx) {
	opts := []cookieOption{WithHTTPOnly(), WithSecure()}
	k.Delete(ctx, opts...)
}

func (k cookieKey) compareValue(ctx *fasthttp.RequestCtx, value any) bool {
	v := string(ctx.Request.Header.Cookie(string(k)))
	return v == fmt.Sprint(value)
}

func (k cookieKey) GetValue(ctx *fasthttp.RequestCtx) string {
	return string(ctx.Request.Header.Cookie(string(k)))
}
