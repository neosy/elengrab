package mappers

import (
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_executor/types"
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/helper"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (m *Mappers) MapDownloadResultToMediaDownload(
	out *ddownload.MediaDownload,
	result *dservices.DownloaderResult,
	thumbnailIDs types.ThumbnailIDs,
) {
	if out == nil || result == nil {
		return
	}

	out.ChannelID = result.ChannelID

	if out.MediaTitle == out.MediaTitleOriginal {
		out.MediaTitle = result.MediaTitle
	}
	out.MediaTitleOriginal = result.MediaTitle

	if helper.ValuesEqual(out.MediaDescription, out.MediaDescriptionOriginal) {
		out.MediaDescription = result.MediaDescription
	}
	out.MediaDescriptionOriginal = result.MediaDescription

	out.Ext = result.FileExt
	out.FileSize = result.Filesize

	out.FileName = result.Filename
	out.FileFullName = result.FileFullName
	out.SafeReadableFileFullName = out.NormalizeFileFullName()

	out.PartialHash = result.PartialHash

	out.MediaInfo = m.MapMediaInfoDomain(result.MediaInfo, thumbnailIDs)
}

func (m *Mappers) MapDownloaderResultToState(out *ddownload.DownloadState, result *dservices.DownloaderResult, mediaInfo *dtypes.MediaInfo) {
	if out == nil {
		return
	}

	if result == nil {
		return
	}

	if result.ChannelID != nil && result.Channel != nil {
		out.Download.ChannelID = result.ChannelID
	}

	if out.Download.MediaTitle == out.Download.MediaTitleOriginal {
		out.Download.MediaTitle = result.MediaTitle
	}
	out.Download.MediaTitleOriginal = result.MediaTitle

	if helper.ValuesEqual(out.Download.MediaDescription, out.Download.MediaDescriptionOriginal) {
		out.Download.MediaDescription = result.MediaDescription
	}
	out.Download.MediaDescriptionOriginal = result.MediaDescription

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
