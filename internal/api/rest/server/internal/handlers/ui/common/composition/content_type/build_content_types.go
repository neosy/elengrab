package contenttype

import (
	"github.com/neosy/elengrab/internal/app/usecases/downloader"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func buildContentTypes(authCtx dauth.AuthContext, downloader downloader.DownloaderAPI, options []contentTypeOption) []dtypes.QueryMediaContentType {
	contentTypes := make([]dtypes.QueryMediaContentType, 0, len(options))

	for _, opt := range options {
		if opt.hasCreateAccess && !downloader.CanCreateMediaDownload(authCtx) {
			continue
		}

		contentTypes = append(contentTypes, opt.contentType)
	}

	return contentTypes
}
