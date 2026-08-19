package types

import (
	dservices "github.com/neosy/elengrab/internal/domain/services"
)

type ProcessedDownload struct {
	Result       *dservices.DownloaderResult
	ThumbnailIDs ThumbnailIDs
}
