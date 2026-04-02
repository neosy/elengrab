package handlers

import (
	"html/template"
	"path/filepath"
	"strings"
	"unicode"

	authmw "github.com/neosy/elengrab/internal/api/rest/server/internal/auth_middleware"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/valyala/fasthttp"
)

type requestFilters map[string]string

func parseFilters(ctx *fasthttp.RequestCtx) requestFilters {
	filters := make(requestFilters)

	for key, value := range ctx.QueryArgs().All() {
		k := string(key)
		v := string(value)

		prefix := "filter["
		suffix := "]"

		if !strings.HasPrefix(k, prefix) || !strings.HasSuffix(k, suffix) {
			continue
		}

		// filter[name] → name
		field := k[len(prefix) : len(k)-len(suffix)]
		if field == "" {
			continue
		}

		filters[field] = v
	}

	return filters
}

func capitalize(s string) string {
	if s == "" {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

func (h *DownloaderHandlers) redirectGuestIfAuthRequired(ctx *fasthttp.RequestCtx) bool {
	if h.appMode == dtypes.AppModeAuthOnly {
		ctxUser := authmw.UserFromContext(ctx)
		if ctxUser == nil || ctxUser.UserType() < dtypes.UserTypeUser {
			ctx.Redirect(httppaths.GroupAccount+httppaths.PathLogin, fasthttp.StatusFound)
			return true
		}
	}
	return false
}

func (h *DownloaderHandlers) loadPage(fileName string) (*template.Template, error) {
	t, err := h.templates.Clone()
	if err != nil {
		return nil, err
	}

	t, err = t.ParseFiles(filepath.Join(h.assetsDir, "templates", "pages", fileName))
	if err != nil {
		return nil, err
	}

	return t, nil
}

func errInternal(err error) error {
	return errorx.Errorf(
		"template execution error: %v", err,
		errorx.ErrorMessageArg("Internal Server Error"),
		errorx.HttpStatusArg(fasthttp.StatusInternalServerError),
	)
}
