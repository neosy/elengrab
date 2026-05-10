package dto

import (
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type FileInfoPatch struct {
	ChannelID            **string
	MediaTitle           *string
	MediaDescription     **string
	FileName             *string
	Ext                  *string
	FullFileName         *string
	FileSize             **int64
	PartialHash          **string
	SafeReadableFullName *string
	MediaInfo            **dtypes.MediaInfo
}

func PatchToFileDomain(patch *FileInfoPatch, file *ddownload.File) {
	if patch == nil || file == nil {
		return
	}

	if patch.ChannelID != nil {
		file.ChannelID = *patch.ChannelID
	}

	if patch.MediaTitle != nil {
		file.MediaTitle = *patch.MediaTitle
	}

	if patch.MediaDescription != nil {
		file.MediaDescription = *patch.MediaDescription
	}

	if patch.FileName != nil {
		file.FileName = *patch.FileName
	}

	if patch.Ext != nil {
		file.Ext = *patch.Ext
	}

	if patch.FullFileName != nil {
		file.FullFileName = *patch.FullFileName
	}

	if patch.FileSize != nil {
		file.FileSize = *patch.FileSize
	}

	if patch.PartialHash != nil {
		file.PartialHash = *patch.PartialHash
	}

	if patch.SafeReadableFullName != nil {
		file.SafeReadableFullName = *patch.SafeReadableFullName
	}

	if patch.MediaInfo != nil {
		file.MediaInfo = *patch.MediaInfo
	}
}
