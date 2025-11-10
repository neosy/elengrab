package mappers

import (
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
)

func (m *Mappers) MapFileDomainToEntity(file *ddownload.File) *edownload.File {
	return &edownload.File{
		FileId:               file.FileId,
		Title:                file.Title,
		FileName:             file.FileName,
		Ext:                  file.Ext,
		FullName:             file.FullName,
		SafeReadableFullName: file.SafeReadableFullName,
	}
}

func (m *Mappers) MapFileEntityToDomain(efile *edownload.File) *ddownload.File {
	return &ddownload.File{
		FileId:               efile.FileId,
		Title:                efile.Title,
		FileName:             efile.FileName,
		Ext:                  efile.Ext,
		FullName:             efile.FullName,
		SafeReadableFullName: efile.SafeReadableFullName,
		CreatedAt:            efile.CreatedAt,
		UpdatedAt:            efile.UpdatedAt,
	}
}
