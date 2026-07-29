package admin

import (
	"bytes"
	"fmt"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/admin/pages/content"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/admin/pages/page"
	"github.com/valyala/fasthttp"
)

func (h *AdminHandlers) renderContent(ctx *fasthttp.RequestCtx, page page.Page) (*bytes.Buffer, error) {
	content, err := content.NewContent(h.mappers, h.usecases.admin, page)
	if err != nil {
		return nil, err
	}

	contentData, err := content.Build(ctx)
	if err != nil {
		return nil, err
	}

	if contentData == nil {
		return nil, fmt.Errorf("buildContentData returned nil for pageTitle=%s", page.Title)
	}

	var buf bytes.Buffer
	err = h.templates.Base.ExecuteTemplate(&buf, page.ContentTemplateName, contentData)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
