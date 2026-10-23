package downloadpreparer

import (
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

type DownloadPlan struct {
	Args         []string
	FileExtNoDot string
	ExtractInfo  *idto.ExtractInfo
	MediaInfo    *pservices.MediaInfo
}
