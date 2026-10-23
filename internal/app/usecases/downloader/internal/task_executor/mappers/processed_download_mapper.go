package mappers

import (
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/helper"
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/task_executor/types"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (m *Mappers) MapProcessedDownloadToMediaDownload(
	out *ddownload.MediaDownload,
	processed *types.ProcessedDownload,
) {
	out.ChannelID = processed.ChannelID

	if out.MediaTitle == out.MediaTitleOriginal {
		out.MediaTitle = processed.MediaTitle
	}
	out.MediaTitleOriginal = processed.MediaTitleOriginal

	if helper.ValuesEqual(out.MediaDescription, out.MediaDescriptionOriginal) {
		out.MediaDescription = processed.MediaDescription
	}
	out.MediaDescriptionOriginal = processed.MediaDescriptionOriginal

	out.FileExtNoDot = processed.FileExtNoDot
	out.FileSize = processed.Filesize

	out.FileNameNoExt = processed.FileNameNoExt
	out.FileName = processed.FileName
	out.SafeReadableFileName = out.NormalizeFileName()

	out.PartialHash = processed.PartialHash

	out.MediaInfo = processed.MediaInfo
}
