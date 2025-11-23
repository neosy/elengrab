package grabberh

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	htmxvalues "github.com/neosy/elengrab/internal/api/rest/server/handlers/htmx/values"
	"github.com/neosy/elengrab/pkg/nfasthttp"
	"github.com/valyala/fasthttp"
)

// IndexHandlers serves the main page (index.html)
func (h *GrabberHandlers) IndexHandler(ctx *fasthttp.RequestCtx) {
	var (
		itemsTmplBuf    bytes.Buffer
		needLoadHistory bool
	)

	resps, err := h.usecases.Downloader.LoadHistory(ctx, time.Now(), 1)
	if err == nil {
		needLoadHistory = len(resps) != 0
	}

	if needLoadHistory {
		itemsTmpl, err := template.ParseFiles(filepath.Join(h.assetsDir, "templates", htmxvalues.GrabResultItemsHistoryHtmlFileName))
		if err != nil {
			nfasthttp.WriteError(ctx, err, fasthttp.StatusInternalServerError)
			return
		}

		if err := itemsTmpl.Execute(&itemsTmplBuf, htmxvalues.PathValues); err != nil {
			nfasthttp.WriteError(ctx, fmt.Errorf("template execution error: %v", err), fasthttp.StatusInternalServerError)
			return
		}
	} else { // Подготовка шаблона с пустой строкой
		itemsTmpl, err := template.ParseFiles(filepath.Join(h.assetsDir, "templates", htmxvalues.GrabResultItemReplacemeHtmlFileName))
		if err != nil {
			nfasthttp.WriteError(ctx, err, fasthttp.StatusInternalServerError)
			return
		}

		dataMap := make(map[string]any)
		dataMap[htmxvalues.DataOnlyOneKey] = "true"

		if err := itemsTmpl.Execute(&itemsTmplBuf, dataMap); err != nil {
			nfasthttp.WriteError(ctx, fmt.Errorf("template execution error: %v", err), fasthttp.StatusInternalServerError)
			return
		}
	}

	tmplPath := filepath.Join(h.assetsDir, "templates", htmxvalues.IndexHtmlFileName)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		nfasthttp.WriteError(ctx, err, fasthttp.StatusInternalServerError)
		return
	}

	// Set content type so browser renders HTML properly
	ctx.SetContentType("text/html; charset=utf-8")

	dataMap := htmxvalues.MergeMaps(htmxvalues.IndexValues, htmxvalues.PathValues)
	dataMap[htmxvalues.GrabResultItemsHtmlKey] = template.HTML(itemsTmplBuf.String())

	// Execute template with PageTitle
	if err := tmpl.Execute(ctx, dataMap); err != nil {
		nfasthttp.WriteError(ctx, fmt.Errorf("template execution error: %v", err), fasthttp.StatusInternalServerError)
		return
	}
}
