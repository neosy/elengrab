package download

import (
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type mediaDownloadQueryOptions struct {
	dtypes.QueryMediaOptions

	includeDeleted bool
	statuses       []dtypes.MediaDownloadStatus

	partialHash **string
}
