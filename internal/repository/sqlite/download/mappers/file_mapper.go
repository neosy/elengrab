package mappers

import (
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
)

func (m *Mappers) MapFileDomainToEntity(file *ddownload.File) *edownload.File {
	return &edownload.File{
		FileId:               file.FileId,
		Status:               file.Status.String(),
		YoutubeTitle:         file.YoutubeTitle,
		FileName:             file.FileName,
		Ext:                  file.Ext,
		FullName:             file.FullName,
		SafeReadableFullName: file.SafeReadableFullName,
		ErrorMessage:         file.ErrorMessage,
	}
}

func (m *Mappers) MapFileEntityToDomain(eFile *edownload.File, eTask *edownload.DownloadTask) *ddownload.File {
	var task *ddownload.DownloadTask
	if eTask != nil && eTask.TaskId.String != "" {
		task = m.MapDownloadTaskEntityToDomain(eTask)
	}

	return &ddownload.File{
		FileId:               eFile.FileId,
		Status:               dtypes.FileStatus(eFile.Status),
		YoutubeTitle:         eFile.YoutubeTitle,
		FileName:             eFile.FileName,
		Ext:                  eFile.Ext,
		FullName:             eFile.FullName,
		SafeReadableFullName: eFile.SafeReadableFullName,
		ErrorMessage:         eFile.ErrorMessage,
		CreatedAt:            eFile.CreatedAt,
		UpdatedAt:            eFile.UpdatedAt,
		DownloadTask:         task,
	}
}
