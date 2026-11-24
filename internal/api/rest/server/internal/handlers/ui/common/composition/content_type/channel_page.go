package contenttype

import (
	"github.com/neosy/elengrab/internal/app/usecases/downloader"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

var channelContentTypeOptions = []contentTypeOption{
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
}

func BuildChannelPageContentTypes(authCtx dauth.AuthContext, downloader downloader.DownloaderAPI) []dtypes.QueryMediaContentType {
	return buildContentTypes(authCtx, downloader, channelContentTypeOptions, nil)
}
