package downloader

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/icons"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/dto"
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

func (h *DownloaderHandlers) buildMediaWatchURL(downloadID uuid.UUID) string {
	return strings.TrimSuffix(h.baseURL, "/") +
		httppaths.BuildPathMediaItemWatch(downloadID)
}

func (h *DownloaderHandlers) getVisibilityResponse(info *ucdto.MediaDownloadInfo) *dto.VisibilityResponse {
	response := &dto.VisibilityResponse{
		Visible: info.ShouldShowVisibility(),
		Value:   info.Visibility.String(),
		Label:   info.Visibility.Label(),
	}

	switch info.Visibility {
	case dtypes.MediaVisibilityPrivate:
		response.Icon = icons.MediaPrivateIcon.FileRaw()

	case dtypes.MediaVisibilityPublic:
		response.Icon = icons.MediaPublicIcon.FileRaw()
	}

	return response
}

func (h *DownloaderHandlers) hasShareLink(ctx context.Context, downloadID uuid.UUID) bool {
	link, _ := h.linkWeb.ResolveURL(ctx, h.buildMediaWatchURL(downloadID))
	return link != nil
}
