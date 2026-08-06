package downloadpreparer

import (
	"context"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
)

type DownloadPreparer struct {
	getExtractInfo func(
		ctx context.Context,
		url string, formatQuery string,
		opts ...idto.ExecutorOption,
	) (*idto.ExtractInfo, error)
}

func NewDownloadPreparer(
	getExtractInfo func(
		ctx context.Context,
		url string, formatQuery string,
		opts ...idto.ExecutorOption,
	) (*idto.ExtractInfo, error),
) *DownloadPreparer {
	return &DownloadPreparer{
		getExtractInfo: getExtractInfo,
	}
}
