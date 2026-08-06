package downloadpreparer

import (
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	dservices "github.com/neosy/elengrab/internal/domain/services"
)

type DownloadPlan struct {
	Args        []string
	FileExt     string
	ExtractInfo *idto.ExtractInfo
	MediaInfo   *dservices.MediaInfo
}
