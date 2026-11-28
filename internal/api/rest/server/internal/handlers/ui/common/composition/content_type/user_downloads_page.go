package contenttype

import (
	"github.com/neosy/elengrab/internal/app/usecases/downloader"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

var userDownloadsContentTypeOptions = []contentTypeOption{
	{
		contentType: dtypes.QueryMediaContentTypeAll,
	},
	{
		contentType: dtypes.QueryMediaContentTypeVideo,
	},
	{
		contentType: dtypes.QueryMediaContentTypeShorts,
	},
	{
		contentType: dtypes.QueryMediaContentTypeAudio,
	},
	{
		contentType: dtypes.QueryMediaContentTypeUploads,
	},
	{
		contentType:     dtypes.QueryMediaContentTypeInProgress,
		hasCreateAccess: true,
	},
}

func BuildUserDownloadsPageContentTypes(
	authCtx dauth.AuthContext,
	downloader downloader.DownloaderAPI,
	isCurrentUserPage bool,
) []dtypes.QueryMediaContentType {
	return buildContentTypes(
		authCtx, downloader, userDownloadsContentTypeOptions,
		func(opt contentTypeOption) bool {
			if opt.contentType == dtypes.QueryMediaContentTypeUploads {
				return isCurrentUserPage || authCtx.IsAdmin()
			}
			if opt.contentType == dtypes.QueryMediaContentTypeInProgress {
				return isCurrentUserPage
			}
			return true
		},
	)
}
