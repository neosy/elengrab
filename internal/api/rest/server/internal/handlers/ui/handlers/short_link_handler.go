package handlers

import (
	"fmt"
	"strings"

	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	"github.com/neosy/elengrab/internal/pkg/nfasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) ShortLinkHandler(ctx *fasthttp.RequestCtx) {
	shortCode := ctx.UserValue(shortCodeKey).(string)
	if shortCode == "" {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTP("shortCode is required", fasthttp.StatusBadRequest))
		return
	}

	shortURL, ipAddress, userAgent, referrer := h.extractRequestMeta(ctx)

	link, err := h.linkWeb.ShortLinkClick(ctx, shortURL, ipAddress, userAgent, referrer)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	if link == nil {
		nfasthttp.WriteErrorx(
			ctx,
			errorx.New(
				"link not found",
				exceptionx.NOT_FOUND,
				exceptionx.NOT_FOUND.ErrorMessage(),
			))
		return
	}

	dataMap := uivalues.MergeMaps(
		uivalues.BaseValues.ToMap(),
		uivalues.PathValues,
	)
	dataMap[uivalues.PathStreamKey] = httppaths.BuildPathStreamShortCode(shortCode)

	// Set content type so browser renders HTML properly
	ctx.SetContentType("text/html; charset=utf-8")

	// Load template
	tmpl, err := h.loadPage(uivalues.PageWatch.FileName())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}

	// Execute template with PageTitle
	if err := tmpl.ExecuteTemplate(ctx, uivalues.PageWatch.Key(), dataMap); err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}
}

// extractRequestMeta extracts metadata from the incoming HTTP request.
// It returns the full request URL, client IP address, user agent, and referrer
func (h *DownloaderHandlers) extractRequestMeta(ctx *fasthttp.RequestCtx) (url string, ip string, userAgent, referrer string) {
	// Определяем схему протокола (http или https)
	scheme := "http"
	if ctx.IsTLS() {
		scheme = "https"
	}
	// Формируем полный URL
	url = fmt.Sprintf("%s://%s%s", scheme, ctx.Host(), ctx.Request.URI().RequestURI())

	// Получаем IP клиента
	ip = ctx.RemoteIP().String()

	// Если заголовок X-Forwarded-For присутствует, берём первый IP из списка
	xff := ctx.Request.Header.Peek("X-Forwarded-For")
	if len(xff) > 0 {
		xffList := strings.Split(string(xff), ",")
		if len(xffList) > 0 {
			xffFirst := strings.TrimSpace(xffList[0])
			if xffFirst != "" {
				ip = xffFirst
			}
		}
	}

	// Получаем User-Agent
	userAgent = string(ctx.Request.Header.UserAgent())

	// Получаем Referer
	referrer = string(ctx.Request.Header.Referer())

	return
}
