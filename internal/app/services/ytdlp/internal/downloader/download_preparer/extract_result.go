package downloadpreparer

import (
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type extractResult struct {
	formatQuery string
	args        []string

	fileExt string

	info        *idto.ExtractInfo
	mediaFormat *idto.ExtractMediaFormat

	videoCodec dtypes.VideoCodec
	audioCodec dtypes.AudioCodec

	selectedFormatIDs []string
}
