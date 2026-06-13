package downloader

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

type cookieKey string

func (k cookieKey) makeCookie(value any, path string, maxAgeSec int) *fasthttp.Cookie {
	var cookie fasthttp.Cookie

	cookie.SetKey(string(k))
	cookie.SetValue(fmt.Sprint(value))
	cookie.SetPath(path)
	cookie.SetMaxAge(maxAgeSec)

	return &cookie
}

func (k cookieKey) compareValue(ctx *fasthttp.RequestCtx, value any) bool {
	v := string(ctx.Request.Header.Cookie(string(k)))
	return v == fmt.Sprint(value)
}
