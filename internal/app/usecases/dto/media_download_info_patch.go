package dto

import (
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type MediaDownloadInfoPatch struct {
	ChannelID            **string
	MediaTitle           *string
	MediaDescription     **string
	FileName             *string
	Ext                  *string
	FileFullName         *string
	FileSize             **int64
	PartialHash          **string
	SafeReadableFullName *string
	MediaInfo            **dtypes.MediaInfo
}

func PatchToMediaDownloadDomain(patch *MediaDownloadInfoPatch, download *ddownload.MediaDownload) {
	if patch == nil || download == nil {
		return
	}

	if patch.ChannelID != nil {
		download.ChannelID = *patch.ChannelID
	}

	if patch.MediaTitle != nil {
		download.MediaTitle = *patch.MediaTitle
	}

	if patch.MediaDescription != nil {
		download.MediaDescription = *patch.MediaDescription
	}

	if patch.FileName != nil {
		download.FileName = *patch.FileName
	}

	if patch.Ext != nil {
		download.Ext = *patch.Ext
	}

	if patch.FileFullName != nil {
		download.FileFullName = *patch.FileFullName
	}

	if patch.FileSize != nil {
		download.FileSize = *patch.FileSize
	}

	if patch.PartialHash != nil {
		download.PartialHash = *patch.PartialHash
	}

	if patch.SafeReadableFullName != nil {
		download.SafeReadableFullName = *patch.SafeReadableFullName
	}

	if patch.MediaInfo != nil {
		download.MediaInfo = *patch.MediaInfo
	}
}
