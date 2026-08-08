package dto

import "github.com/google/uuid"

type MediaDownloadChanged struct {
	DownloadID uuid.UUID

	TitleChanged      bool
	VisibilityChanged bool
	ShareLinkChanged  bool

	WatchPositionChanged bool
	WatchedChanged       bool
	ViewCountChanged     bool

	Info *MediaDownloadInfo
}

func (m *MediaDownloadChanged) MarkWatchStatsChanged() {
	m.WatchPositionChanged = true
	m.WatchedChanged = true
	m.ViewCountChanged = true
}

func (m *MediaDownloadChanged) MarkWatchPositionChanged() {
	m.WatchPositionChanged = true
}

func (m *MediaDownloadChanged) MarkManualChanges() {
	m.TitleChanged = true
	m.VisibilityChanged = true
}

func (m *MediaDownloadChanged) MarkShareLinkChanges() {
	m.ShareLinkChanged = true
}
