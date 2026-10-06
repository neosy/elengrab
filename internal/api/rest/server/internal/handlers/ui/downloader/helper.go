package downloader

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/uuid"
	apierrors "github.com/neosy/elengrab/internal/api/errors"
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
	dlink "github.com/neosy/elengrab/internal/domain/link"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/neosy/elengrab/internal/pkg/stringx"
	"github.com/valyala/fasthttp"
)

type shortLink struct {
	code string

	shortURL   string
	pathSuffix string

	originalURL string
}

func parseGetFilters(ctx *fasthttp.RequestCtx) (*dtypes.QueryFilters, error) {
	filters := dtypes.NewQueryFilters()

	for key, value := range ctx.QueryArgs().All() {
		k := string(key)
		v := string(value)

		switch k {
		case qkeys.SearchKey.String():
			fallthrough
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

func parseSearchQueryString(queryString string) (*types.QueryFilters, error) {
	if queryString == "" {
		return nil, nil
	}

	items := strings.Split(queryString, "&")
	if len(items) == 0 {
		return nil, nil
	}

	queryItems := types.NewQueryFilters()

	for _, item := range items {
		parts := strings.SplitN(item, "=", 2)
		if len(parts) != 2 {
			continue
		}

		value, err := url.QueryUnescape(parts[1])
		if err != nil {
			return nil, err
		}

		k := parts[0]

		key := qkeys.SearchQueryKeys.FindByShortKey(k)
		if key == "" {
			key = qkeys.SearchQueryKeys.FindByStringKey(k)
		}

		if key == "" {
			continue
		}

		if !qkeys.SearchQueryKeys.ExistsByKey(key) {
			continue
		}

		queryItems.Add(key, value)
	}

	if queryItems.Len() == 0 {
		return nil, nil
	}

	return queryItems, nil
}

func (h *DownloaderHandlers) parseSearchQueryString(queryString string) (types.SearchValues, error) {
	if queryString == "" {
		return types.NewSearchValues(), nil
	}

	queryItems, err := parseSearchQueryString(queryString)
	if err != nil {
		return types.SearchValues{}, err
	}

	if queryItems.Len() == 0 {
		return types.NewSearchValues(), nil
	}

	searchValues, err := h.mappers.MapSearchQueryToSearchValues(queryItems)
	if err != nil {
		return types.SearchValues{}, err
	}

	return searchValues, nil
}

func (h *DownloaderHandlers) parseSearchGetRequest(ctx *fasthttp.RequestCtx) (types.SearchValues, error) {
	return h.parseSearchQueryString(ctx.QueryArgs().String())
}

func (h *DownloaderHandlers) parseSearchPostRequest(ctx *fasthttp.RequestCtx) (types.SearchValues, error) {
	return h.parseSearchQueryString(ctx.PostArgs().String())
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
		httppaths.BuildMediaItemPath(downloadID)
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

	channelURL := httppaths.BuildChannelPath(channelID)

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
			Username: channel.Username,
		},
		Show: true,
	}
}

func (h *DownloaderHandlers) extractDownloadID(ctx *fasthttp.RequestCtx) (uuid.UUID, error) {
	downloadIDStr, ok := ctx.UserValue(qkeys.DownloadIDKey.String()).(string)
	if !ok || downloadIDStr == "" {
		return uuid.Nil, apierrors.ErrDownloadIDIsRequired
	}

	downloadID, err := idcodec.DecodeUUIDBase64URL(downloadIDStr)
	if err != nil {
		return uuid.Nil, apierrors.ErrDownloadIDIsIncorrect.Wrap(err)
	}

	return downloadID, nil
}

func (h *DownloaderHandlers) extractShortURL(rawURL string) (string, string) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", ""
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 || "/"+parts[0] != h.shortLinkPrefix {
		return "", ""
	}

	shortURL := u.Scheme + "://" + u.Host + "/" + strings.Join(parts[:2], "/")
	pathSuffix := strings.Join(parts[2:], "/")

	return shortURL, pathSuffix
}

func (h *DownloaderHandlers) resolveShortLink(ctx *fasthttp.RequestCtx, clickShortLink bool) (*shortLink, error) {
	shortCode, ok := ctx.UserValue(qkeys.ShortCodeKey.String()).(string)
	if !ok || shortCode == "" {
		return nil, errorx.NewHTTPMessage("shortCode is required", fasthttp.StatusBadRequest)
	}

	url, ipAddress, userAgent, referrer := h.extractRequestMeta(ctx)

	shortURL, pathSuffix := h.extractShortURL(url)
	if shortURL == "" {
		return nil, errorx.NewHTTPMessage("short link URL is invalid", fasthttp.StatusBadRequest)
	}

	var (
		link *dlink.Link
		err  error
	)

	if clickShortLink {
		link, err = h.linkWeb.ShortLinkClick(ctx, shortURL, ipAddress, userAgent, referrer)
	} else {
		link, err = h.linkWeb.GetLastByShortCode(ctx, shortCode)
	}

	if err != nil {
		return nil, err
	}

	if link == nil {
		return nil, errorx.NewMessage("short link not found.", exceptionx.NOT_FOUND)
	}

	return &shortLink{
		code: shortCode,

		shortURL:   shortURL,
		pathSuffix: pathSuffix,

		originalURL: link.OriginalURL,
	}, nil
}

func (h *DownloaderHandlers) resolveShortLinkToDownloadID(ctx *fasthttp.RequestCtx, clickShortLink bool) (*shortLink, uuid.UUID, error) {
	shortLink, err := h.resolveShortLink(ctx, clickShortLink)
	if err != nil {
		return nil, uuid.Nil, err
	}

	downloadID := stripUUIDFromIDPath(shortLink.originalURL)
	if downloadID == uuid.Nil {
		return nil, uuid.Nil, errorx.Errorf(
			"failed short link %v, originalURL: %v", shortLink.shortURL, shortLink.originalURL,
			exceptionx.WRONG_DATA,
			errorx.WithErrorMessage("Short link is invalid."),
		)
	}

	return shortLink, downloadID, nil
}
