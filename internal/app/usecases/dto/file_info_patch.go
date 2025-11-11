package dto

import ddownload "github.com/neosy/elengrab/internal/domain/download"

type FileInfoPatch struct {
	YoutubeTitle         *string
	FileName             *string
	Ext                  *string
	FullName             *string
	SafeReadableFullName *string
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

	if patch.SafeReadableFullName != nil {
		file.SafeReadableFullName = *patch.SafeReadableFullName
	}
}
