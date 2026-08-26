package mappers

import (
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_executor/types"
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/helper"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (m *Mappers) MapDownloadResultToProcessedDownload(
	result *dservices.DownloaderResult,
	state *ddownload.DownloadState,
	thumbnailIDs types.ThumbnailIDs,
) *types.ProcessedDownload {
	title := result.MediaTitle
	description := result.MediaDescription

	if state != nil && state.Download != nil {
		title = state.Download.MediaTitle
		description = state.Download.MediaDescription
	}

	var mediaInfo *dtypes.MediaInfo
	if result.MediaInfo != nil {
		mediaInfo = m.MapMediaInfoDomain(result.MediaInfo, thumbnailIDs)
	}

	return &types.ProcessedDownload{
		MediaTitle:         title,
		MediaTitleOriginal: result.MediaTitle,

		MediaDescription:         description,
		MediaDescriptionOriginal: result.MediaDescription,

		Filename:     result.Filename,
		FileFullName: result.FileFullName,

		FileExt:  result.FileExt,
		Filesize: result.Filesize,

		PartialHash: result.PartialHash,

		ChannelID: result.ChannelID,

		MediaInfo: mediaInfo,
	}
}

func (m *Mappers) MapDownloaderResultToState(
	out *ddownload.DownloadState,
	result *dservices.DownloaderResult,
	mediaInfo *dtypes.MediaInfo,
) {
	if out.Download.MediaTitle == out.Download.MediaTitleOriginal {
		out.Download.MediaTitle = result.MediaTitle
	}
	out.Download.MediaTitleOriginal = result.MediaTitle

	if helper.ValuesEqual(out.Download.MediaDescription, out.Download.MediaDescriptionOriginal) {
		out.Download.MediaDescription = result.MediaDescription
	}
	out.Download.MediaDescriptionOriginal = result.MediaDescription

	if result.ChannelID != nil && result.Channel != nil {
		out.Download.ChannelID = result.ChannelID
	}

	out.Download.Ext = result.FileExt

	if result.Filesize != nil {
		out.Download.FileSize = result.Filesize
	}

	if result.PartialHash != nil {
		out.Download.PartialHash = result.PartialHash
	}

	if mediaInfo != nil {
		out.Download.MediaInfo = mediaInfo
	}

	if result.Progress != nil {
		out.Progress = result.Progress
	}
}
