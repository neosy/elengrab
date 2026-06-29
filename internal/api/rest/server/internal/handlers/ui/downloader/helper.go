package downloader

import (
	"fmt"
	"html/template"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/neosy/elengrab/internal/pkg/stringx"
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

func (h *DownloaderHandlers) redirectGuestIfAuthRequired(ctx *fasthttp.RequestCtx) bool {
	if h.appMode == dtypes.AppModeAuthenticated {
		ctxUser := policy.ResolveUser(ctx)
		if ctxUser == nil || ctxUser.UserType() < dtypes.UserTypeUser {
			ctx.Redirect(httppaths.GroupAccount+httppaths.PathLogin, fasthttp.StatusFound)
			return true
		}
	}
	return false
}

func (h *DownloaderHandlers) loadPageTemplate(fileName string) (*template.Template, error) {
	t, err := h.templates.Clone()
	if err != nil {
		return nil, err
	}

	t, err = t.ParseFiles(filepath.Join(h.assets.FolderPaths().Pages(), fileName))
	if err != nil {
		return nil, err
	}

	return t, nil
}

func errInternal(err error) error {
	return errorx.Errorf(
		"template execution error: %v", err,
		errorx.NewFromDomainException(exceptionx.ERROR),
	)
}

func (h *DownloaderHandlers) getScheme(ctx *fasthttp.RequestCtx) string {
	if h.baseURL != "" {
		if s := httpx.SchemeFromURL(h.baseURL); s != "" {
			return s
		}
	}

	if proto := string(ctx.Request.Header.Peek("X-Forwarded-Proto")); proto != "" {
		return proto
	}
	if proto := string(ctx.Request.Header.Peek("X-Forwarded-Protocol")); proto != "" {
		return proto
	}
	if ctx.IsTLS() {
		return "https"
	}

	return "http"
}

// extractRequestMeta extracts metadata from the incoming HTTP request.
// It returns the full request URL, client IP address, user agent, and referrer
func (h *DownloaderHandlers) extractRequestMeta(ctx *fasthttp.RequestCtx) (url string, ip string, userAgent, referrer string) {
	// Determine the protocol scheme (http or https)
	scheme := h.getScheme(ctx)

	// Creating the full URL
	url = fmt.Sprintf("%s://%s%s", scheme, ctx.Host(), ctx.Request.URI().RequestURI())

	// Getting the client's IP address
	ip = nfasthttp.GetClientIP(ctx)

	// Getting the User-Agent
	userAgent = string(ctx.Request.Header.UserAgent())

	// Getting the Referer
	referrer = string(ctx.Request.Header.Referer())

	return
}

func stripUUIDFromIDPath(path string) uuid.UUID {
	parts := strings.Split(path, "/")

	for _, p := range parts {
		if p == "" {
			continue
		}

		u, err := idcodec.DecodeUUIDBase64URL(p)
		if err == nil {
			return u
		}
	}

	return uuid.Nil
}

func mediaSourceFromURL(mediaURL string) string {
	source := hostdetect.Detect(mediaURL).Title()
	if source != "" {
		return source
	}

	u, err := url.Parse(mediaURL)
	if err != nil {
		return mediaURL
	}

	source = stringx.Capitalize(u.Hostname())
	if source != "" {
		return source
	}

	return mediaURL
}
