package indexh

import (
	"fmt"
	"html/template"
	"path/filepath"

	avalues "github.com/neosy/elengrab/internal/api/rest/server/assets/values"
	"github.com/neosy/elengrab/pkg/nfasthttp"
	"github.com/valyala/fasthttp"
)

// IndexHandlers serves the main page (index.html)
func (h *IndexHandlers) IndexHandler(ctx *fasthttp.RequestCtx) {
	tmplPath := filepath.Join(h.assetsDir, "templates", "index.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		nfasthttp.WriteError(ctx, err, fasthttp.StatusInternalServerError)
		return
	}

	// Set content type so browser renders HTML properly
	ctx.SetContentType("text/html; charset=utf-8")

	// Execute template with PageTitle
	if err := tmpl.Execute(ctx, avalues.MergeMaps(avalues.IndexValues, avalues.PathValues)); err != nil {
		nfasthttp.WriteError(ctx, fmt.Errorf("template execution error: %v", err), fasthttp.StatusInternalServerError)
		return
	}
}
