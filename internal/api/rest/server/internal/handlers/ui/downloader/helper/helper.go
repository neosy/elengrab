package helper

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash/fnv"
	"net/url"
	"strings"
	"unicode/utf8"

	"uuid"

	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/icons"
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
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/neosy/elengrab/internal/pkg/stringx"
	"github.com/valyala/fasthttp"
)

var baseURL string

// Init initializes the helper package.
func Init(url string) {
	baseURL = url
}

func ParseSearchQueryFilters(queryString string) (*types.QueryFilters, error) {
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

func ParseSearchQueryValues(queryString string) (types.SearchValues, error) {
	if queryString == "" {
		return types.NewSearchValues(), nil
	}

	queryItems, err := ParseSearchQueryFilters(queryString)
	if err != nil {
		return types.SearchValues{}, err
	}

	if queryItems.Len() == 0 {
		return types.NewSearchValues(), nil
	}

	searchValues, err := mapSearchQueryToSearchValues(queryItems)
	if err != nil {
		return types.SearchValues{}, err
	}

	return searchValues, nil
}

func ParseSearchGetRequest(ctx *fasthttp.RequestCtx) (types.SearchValues, error) {
	return ParseSearchQueryValues(ctx.QueryArgs().String())
}

func ParseSearchPostRequest(ctx *fasthttp.RequestCtx) (types.SearchValues, error) {
	return ParseSearchQueryValues(ctx.PostArgs().String())
}

func ErrTemplateExecution(err error) error {
	return errorx.Errorf(
		"template execution error: %v", err,
		errorx.NewFromDomainException(exceptionx.ERROR),
	)
}

func StripUUIDFromIDPath(path string) uuid.UUID {
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

	return uuid.Nil()
}

func MediaSourceFromURL(mediaURL string) string {
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

func GetScheme(ctx *fasthttp.RequestCtx) string {
	if baseURL != "" {
		if s := httpx.SchemeFromURL(baseURL); s != "" {
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

func GetBaseURL(ctx *fasthttp.RequestCtx) string {
	if baseURL != "" {
		return baseURL
	}

	if ctx == nil || len(ctx.Host()) == 0 {
		return ""
	}

	return GetScheme(ctx) + "://" + string(ctx.Host())
}

// extractRequestMeta extracts metadata from the incoming HTTP request.
// It returns the full request URL, client IP address, user agent, and referrer
func ExtractRequestMeta(ctx *fasthttp.RequestCtx) (url string, ip string, userAgent, referrer string) {
	// Determine the protocol scheme (http or https)
	scheme := GetScheme(ctx)

	// Creating the full URL
	url = fmt.Sprintf("%s://%s%s", scheme, ctx.Host(), ctx.Request.URI().RequestURI())

	// Getting the client's IP address
	ip = fasthttpx.GetClientIP(ctx)

	// Getting the User-Agent
	userAgent = string(ctx.Request.Header.UserAgent())

	// Getting the Referer
	referrer = string(ctx.Request.Header.Referer())

	return
}

func EnsureAbsoluteURL(ctx *fasthttp.RequestCtx, rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse URL %q: %w", rawURL, err)
	}

	if parsedURL.IsAbs() {
		return rawURL, nil
	}

	baseURL := GetBaseURL(ctx)
	if baseURL == "" {
		return "", errors.New("base URL is empty")
	}

	return strings.TrimSuffix(baseURL, "/") + "/" + strings.TrimPrefix(rawURL, "/"), nil
}

func BuildMediaWatchURL(ctx *fasthttp.RequestCtx, downloadID uuid.UUID) string {
	return strings.TrimSuffix(GetBaseURL(ctx), "/") +
		httppaths.BuildMediaItemPath(downloadID)
}

func shouldShowVisibility(visibility dtypes.MediaVisibility) bool {
	_, exists := consts.DisplayedVisibilities[visibility]
	return exists
}

func BuildVisibilityResponse(info *ucdto.MediaDownloadInfo) *dto.VisibilityResponse {
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

func BuildMediaAvatarImageURL(downloadInfo *ucdto.MediaDownloadInfo) string {
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

func BuildMediaSiteImageURL(downloadInfo *ucdto.MediaDownloadInfo) string {
	url := httppaths.BuildMediaItemImagePath(
		downloadInfo.DownloadID,
		downloadInfo.ImageMetaHash(),
		[]dtypes.ImageSource{
			dtypes.ImageSourceSite,
		},
	)
	return url
}

func BuildUserDownloadsHeaderPageData(
	ctx context.Context,
	login, userName string,
	isCurrentUserPage bool,
) pages.UserDownloadsHeader {
	if isCurrentUserPage {
		userName = "My Downloads"
	}

	return pages.UserDownloadsHeader{
		Title:       userName,
		Description: login,
		ImageURL:    "",
		Show:        true,
	}
}

func ExtractDownloadID(ctx *fasthttp.RequestCtx) (uuid.UUID, error) {
	encodedDownloadID, ok := ctx.UserValue(qkeys.DownloadIDKey.String()).(string)
	if !ok || encodedDownloadID == "" {
		return uuid.Nil(), apierrors.ErrDownloadIDIsRequired
	}

	downloadID, err := idcodec.DecodeUUIDBase64URL(encodedDownloadID)
	if err != nil {
		return uuid.Nil(), apierrors.ErrDownloadIDIsIncorrect.Wrap(err)
	}

	return downloadID, nil
}

func ExtractFileID(ctx *fasthttp.RequestCtx) (uuid.UUID, error) {
	encodedFileID, ok := ctx.UserValue(qkeys.FileIDKey.String()).(string)
	if !ok || encodedFileID == "" {
		return uuid.Nil(), apierrors.ErrFileIDIsRequired
	}

	fileID, err := idcodec.DecodeUUIDBase64URL(encodedFileID)
	if err != nil {
		return uuid.Nil(), apierrors.ErrFileIDIsIncorrect.Wrap(err)
	}

	return fileID, nil
}

func ExtractChannelID(ctx *fasthttp.RequestCtx) (uuid.UUID, error) {
	encodedChannelID, ok := ctx.UserValue(qkeys.ChannelIDKey.String()).(string)
	if !ok || encodedChannelID == "" {
		return uuid.Nil(), apierrors.ErrChannelIDIsRequired
	}

	channelID, err := idcodec.DecodeUUIDBase64URL(encodedChannelID)
	if err != nil {
		return uuid.Nil(), apierrors.ErrChannelIDIsIncorrect.Wrap(err)
	}

	return channelID, nil
}

func UserDisplayName(userID uuid.UUID, userLogin string, userType dtypes.UserType, curUserID uuid.UUID) string {
	switch userType {
	case dtypes.UserTypeGuest:
		if curUserID == userID {
			return "Guest (You)"
		}

		id := strings.TrimPrefix(userLogin, "u-")
		num := generateGuestNumber(id)

		return fmt.Sprintf("Guest #%04d", num)
	}

	return stringx.Capitalize(userLogin)
}

// UserDisplayInitial returns the uppercased first letter of the name
func UserDisplayInitial(userName string) string {
	r, _ := utf8.DecodeRuneInString(strings.TrimSpace(userName))
	if r == utf8.RuneError {
		return "?"
	}
	return strings.ToUpper(string(r))
}

// UserDisplayInitialHueFNV derives a stable color hue (0-359) from the user ID using FNV-1a
func UserDisplayInitialHueFNV(userID string) int {
	if userID == "" || userID == nilUUIDString {
		return defaultAvatarHue
	}

	h := fnv.New32a()
	h.Write([]byte(userID))

	return int(h.Sum32() % 360)
}

// UserDisplayInitialHueSHA256 derives a stable color hue (0-359) from the user ID using SHA-256
func UserDisplayInitialHueSHA256(userID string) int {
	if userID == "" || userID == nilUUIDString {
		return defaultAvatarHue
	}

	h := sha256.Sum256([]byte(userID))

	return int(uint16(h[0])<<8|uint16(h[1])) % 360
}

func generateGuestNumber(id string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(id))

	return h.Sum32() % 10000
}
