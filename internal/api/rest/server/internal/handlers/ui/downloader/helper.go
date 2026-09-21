package downloader

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/icons"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/consts"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/types"
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

func parseGetFilters(ctx *fasthttp.RequestCtx) (*dtypes.QueryFilters, error) {
	filters := dtypes.NewQueryFilters()

	for key, value := range ctx.QueryArgs().All() {
		k := string(key)
		v := string(value)

		switch k {
		case qkeys.SearchQueryKey.String():
			filters.Add(dtypes.QueryFilterNameSearchQuery, v)
			continue
		}

		prefix := "filter["
		suffix := "]"

		if !strings.HasPrefix(k, prefix) || !strings.HasSuffix(k, suffix) {
			continue
		}

		// filter[name] → name
		name := k[len(prefix) : len(k)-len(suffix)]
		if name == "" {
			continue
		}

		filterName, err := dtypes.ParseQueryFilterName(name)
		if err != nil {
			return nil, err
		}

		filters.Add(filterName, v)

	}

	return filters, nil
}

func parseGetSearchParameters(ctx *fasthttp.RequestCtx) (*types.SearchParameters, error) {
	searchParametersStr := string(ctx.QueryArgs().Peek(qkeys.SearchParametersKey.Short()))
	if searchParametersStr == "" {
		searchParametersStr = string(ctx.QueryArgs().Peek(qkeys.SearchParametersKey.String()))
	}

	if searchParametersStr == "" {
		return nil, nil
	}

	searchParameters, err := types.ParseEncodedSearchParamQueryString(searchParametersStr)
	if err != nil {
		return nil, err
	}

	return searchParameters, nil
}

func parsePostSearchParameters(ctx *fasthttp.RequestCtx) (*types.SearchParameters, error) {
	queryString := ctx.PostArgs().String()

	if queryString == "" {
		return nil, nil
	}

	searchParameters, err := types.ParseSearchQueryString(queryString)
	if err != nil {
		return nil, err
	}

	return searchParameters, nil
}

func (h *DownloaderHandlers) redirectGuestIfAuthRequired(ctx *fasthttp.RequestCtx) bool {
	if h.appMode == dtypes.AppModeAuthenticated {
		ctxUser := policy.ResolveUser(ctx)
		if ctxUser == nil || ctxUser.UserType() < dtypes.UserTypeUser {
			ctx.Redirect(httppaths.AuthLoginPath, fasthttp.StatusFound)
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
	source := hostdetect.DetectPlatformType(mediaURL).Title()
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
		httppaths.BuildMediaItemWatchPath(downloadID)
}

func shouldShowVisibility(visibility dtypes.MediaVisibility) bool {
	_, exists := consts.DisplayedVisibilities[visibility]
	return exists
}

func (h *DownloaderHandlers) buildVisibilityResponse(info *ucdto.MediaDownloadInfo) *dto.VisibilityResponse {
	response := &dto.VisibilityResponse{
		Visible: shouldShowVisibility(info.Visibility),
		Value:   info.Visibility.String(),
		Label:   info.Visibility.Label(),
	}

	switch info.Visibility {
	case dtypes.MediaVisibilityPrivate:
		response.Icon = icons.MediaPrivateIcon.FileRaw()
	case dtypes.MediaVisibilityAuthenticated:
		response.Icon = icons.MediaAuthenticatedIcon.FileRaw()
	case dtypes.MediaVisibilityPublic:
		response.Icon = icons.MediaPublicIcon.FileRaw()
	}

	return response
}

func (h *DownloaderHandlers) hasShareLink(ctx context.Context, downloadID uuid.UUID) bool {
	link, _ := h.linkWeb.ResolveURL(ctx, h.buildMediaWatchURL(downloadID))
	return link != nil
}

func (h *DownloaderHandlers) buildMediaAvatarImageURL(downloadInfo *ucdto.MediaDownloadInfo) string {
	url := httppaths.BuildMediaItemImagePath(
		downloadInfo.DownloadID,
		downloadInfo.ImageMetaHash(),
		[]dtypes.ImageSource{
			dtypes.ImageSourceChannel,
			dtypes.ImageSourceSite,
		},
	)
	return url
}

func (h *DownloaderHandlers) buildMediaSiteImageURL(downloadInfo *ucdto.MediaDownloadInfo) string {
	url := httppaths.BuildMediaItemImagePath(
		downloadInfo.DownloadID,
		downloadInfo.ImageMetaHash(),
		[]dtypes.ImageSource{
			dtypes.ImageSourceSite,
		},
	)
	return url
}

func (h *DownloaderHandlers) buildChannelPageData(
	ctx context.Context,
	channelID uuid.UUID,
) pages.Channel {
	if channelID == uuid.Nil {
		return pages.Channel{}
	}

	channel, _ := h.downloader.GetChannelByID(ctx, channelID)
	if channel == nil {
		return pages.Channel{}
	}

	encodeChannelID := idcodec.EncodeUUIDBase64URL(channelID)

	searchParameters := types.NewSearchParameters()
	searchParameters.Add(qkeys.ViewModeKey, dtypes.QueryMediaViewModeDefault.String())
	searchParameters.Add(qkeys.ChannelIDKey, encodeChannelID)

	channelURL := "/?" + searchParameters.EncodeShortQueryString()

	return pages.Channel{
		EncodedChannelID: encodeChannelID,
		Title:            channel.Title,

		URL:      channelURL,
		ImageURL: httppaths.BuildChannelImagePath(channel.ExternalID, channel.Platform),

		Ext: pages.ChannelExt{
			ExtID:    channel.ExternalID,
			Platform: channel.Platform,
			URL:      channel.ChannelURL,
		},
	}
}

func (h *DownloaderHandlers) buildChannelHeaderPageData(ctx context.Context, channelID uuid.UUID) pages.ChannelHeader {
	if channelID == uuid.Nil {
		return pages.ChannelHeader{}
	}

	channel, _ := h.downloader.GetChannelByID(ctx, channelID)
	if channel == nil {
		return pages.ChannelHeader{}
	}

	return pages.ChannelHeader{
		EncodedChannelID: idcodec.EncodeUUIDBase64URL(channelID),
		Title:            channel.Title,

		ImageURL: httppaths.BuildChannelImagePath(channel.ExternalID, channel.Platform),

		Ext: pages.ChannelExt{
			ExtID:    channel.ExternalID,
			Platform: channel.Platform,
			URL:      channel.ChannelURL,
		},
		Show: true,
	}
}
