package ddownload

import (
	"github.com/google/uuid"
)

type DownloadState struct {
	FileId uuid.UUID
	TaskId *uuid.UUID

	File *File
}

func (s *DownloadState) InitFromFile(file *File) {
	if s == nil {
		return
	}

	if file == nil {
		return
	}

	s.FileId = file.FileId
	s.File = file

	if file.DownloadTask != nil {
		s.TaskId = &file.DownloadTask.TaskId
	}
}

func (s *DownloadState) InitFromDownloadResult(result *DownloadResult) {
	if s == nil {
		return
	}

	if result == nil {
		return
	}

	s.File.YoutubeChannelID = result.ChannelID
	s.File.YoutubeTitle = result.YoutubeTitle
	s.File.Ext = result.FileExt
	s.File.FileSize = result.Filesize
	s.File.PartialHash = result.PartialHash
}
