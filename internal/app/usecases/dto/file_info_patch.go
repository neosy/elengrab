package dto

import ddownload "github.com/neosy/elengrab/internal/domain/download"

type FileInfoPatch struct {
	YoutubeTitle         *string
	FileName             *string
	Ext                  *string
	FullName             *string
	FileSize             **int
	PartialHash          **string
	SafeReadableFullName *string
	MediaInfo            **ddownload.MediaInfo
}

func PatchToFileDomain(patch *FileInfoPatch, file *ddownload.File) {
	if patch.YoutubeTitle != nil {
		file.YoutubeTitle = *patch.YoutubeTitle
	}

	if patch.FileName != nil {
		file.FileName = *patch.FileName
	}

	if patch.Ext != nil {
		file.Ext = *patch.Ext
	}

	if patch.FullName != nil {
		file.FullName = *patch.FullName
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
