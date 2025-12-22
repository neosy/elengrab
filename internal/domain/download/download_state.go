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

	if result.ChannelID != nil {
		s.File.YoutubeChannelID = result.ChannelID
	}
	if result.YoutubeTitle != nil {
		s.File.YoutubeTitle = *result.YoutubeTitle
	}
	if result.FileExt != nil {
		s.File.Ext = *result.FileExt
	}
	if result.Filesize != nil {
		s.File.FileSize = result.Filesize
	}
	if result.PartialHash != nil {
		s.File.PartialHash = result.PartialHash
	}
}
