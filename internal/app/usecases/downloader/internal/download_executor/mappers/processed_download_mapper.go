package mappers

import (
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_executor/types"
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/helper"
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

	out.Ext = processed.FileExt
	out.FileSize = processed.Filesize

	out.FileName = processed.Filename
	out.FileFullName = processed.FileFullName
	out.SafeReadableFileFullName = out.NormalizeFileFullName()

	out.PartialHash = processed.PartialHash

	out.MediaInfo = processed.MediaInfo
}
