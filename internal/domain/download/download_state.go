package ddownload

import (
	"github.com/google/uuid"
)

type DownloadState struct {
	UserID *uuid.UUID
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

	s.UserID = file.UserID
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

	if result.ChannelID != nil && result.ChannelAvatar != nil {
		s.File.YoutubeChannelID = result.ChannelID
	}

	s.File.YoutubeTitle = result.YoutubeTitle
	s.File.Ext = result.FileExt

	if result.Filesize != nil {
		s.File.FileSize = result.Filesize
	}

	if result.PartialHash != nil {
		s.File.PartialHash = result.PartialHash
	}
}
