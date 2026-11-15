package downloader

import (
	"context"
	"net/url"
	"strings"
	"uuid"

	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/helper"
	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/types"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	dlink "github.com/neosy/elengrab/internal/domain/link"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/valyala/fasthttp"
)

type shortLink struct {
	code string

	shortURL   string
	pathSuffix string

	originalURL string
}

func (h *DownloaderHandlers) redirectIfAuthRequired(ctx *fasthttp.RequestCtx) bool {
	if h.appMode == dtypes.AppModeAuthenticated {
		ctxUser := policy.ResolveUser(ctx)
		if ctxUser == nil || ctxUser.UserType() < dtypes.UserTypeUser {
			ctx.Redirect(httppaths.AuthLoginPath, fasthttp.StatusFound)
			return true
		}
	}
	return false
}

func (h *DownloaderHandlers) hasShareLink(ctx *fasthttp.RequestCtx, downloadID uuid.UUID) bool {
	link, _ := h.linkWeb.ResolveURL(ctx, helper.BuildMediaWatchURL(ctx, downloadID))
	return link != nil
}

func (h *DownloaderHandlers) buildChannelPageData(
	ctx context.Context,
	channelID uuid.UUID,
) pages.Channel {
	if channelID == uuid.Nil() {
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
	if channelID == uuid.Nil() {
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

func (h *DownloaderHandlers) resolveShortLink(ctx *fasthttp.RequestCtx, clickShortLink bool) (*shortLink, error) {
	shortCode, ok := ctx.UserValue(qkeys.ShortCodeKey.String()).(string)
	if !ok || shortCode == "" {
		return nil, errorx.NewHTTPMessage("shortCode is required", fasthttp.StatusBadRequest)
	}

	url, ipAddress, userAgent, referrer := helper.ExtractRequestMeta(ctx)

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
		return nil, uuid.Nil(), err
	}

	downloadID := helper.StripUUIDFromIDPath(shortLink.originalURL)
	if downloadID == uuid.Nil() {
		return nil, uuid.Nil(), errorx.Errorf(
			"failed short link %v, originalURL: %v", shortLink.shortURL, shortLink.originalURL,
			exceptionx.WRONG_DATA,
			errorx.WithErrorMessage("Short link is invalid."),
		)
	}

	return shortLink, downloadID, nil
}

func (h *DownloaderHandlers) prepareSearchQueryFilters(ctx context.Context, queryFilters *types.QueryFilters) (*types.QueryFilters, error) {
	if queryFilters == nil {
		return queryFilters, nil
	}

	searchQueryFilters := queryFilters.Clone()

	if filter, exists := searchQueryFilters.Find(qkeys.UserLoginKey); exists {
		searchQueryFilters.Delete(qkeys.UserLoginKey)

		user, err := h.authWeb.GetByLogin(ctx, filter.Value)
		if err != nil {
			return nil, apierrors.ErrUserNameIsIncorrect.Wrap(err)
		}

		searchQueryFilters.Add(qkeys.UserNameKey, user.Login)
	}

	if filter, exists := searchQueryFilters.Find(qkeys.UserIDKey); exists {
		userID, err := idcodec.DecodeUUIDBase64URL(filter.Value)
		if err != nil {
			return nil, apierrors.ErrUserNameIsIncorrect.Wrap(err)
		}

		searchQueryFilters.Delete(qkeys.UserIDKey)

		user, err := h.authWeb.GetByUserID(ctx, userID)
		if err != nil {
			return nil, apierrors.ErrUserNameIsIncorrect.Wrap(err)
		}

		searchQueryFilters.Add(qkeys.UserNameKey, user.Login)
	}

	return searchQueryFilters, nil
}

func (h *DownloaderHandlers) restoreSearchQueryFilters(ctx context.Context, searchQueryFilters *types.QueryFilters) (*types.QueryFilters, error) {
	if searchQueryFilters == nil {
		return searchQueryFilters, nil
	}

	queryFilters := searchQueryFilters.Clone()

	if filter, exists := queryFilters.Find(qkeys.UserNameKey); exists {
		queryFilters.Delete(filter.Key)

		user, err := h.authWeb.GetByLogin(ctx, filter.Value)
		if err != nil {
			return nil, apierrors.ErrUserNameIsIncorrect.Wrap(err)
		}

		queryFilters.Add(qkeys.UserIDKey, idcodec.EncodeUUIDBase64URL(user.UserID))
	}

	if filter, exists := queryFilters.Find(qkeys.UserLoginKey); exists {
		queryFilters.Delete(filter.Key)

		user, err := h.authWeb.GetByLogin(ctx, filter.Value)
		if err != nil {
			return nil, apierrors.ErrUserNameIsIncorrect.Wrap(err)
		}

		queryFilters.Add(qkeys.UserIDKey, idcodec.EncodeUUIDBase64URL(user.UserID))
	}

	return queryFilters, nil
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
