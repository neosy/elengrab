package dto

import "uuid"

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

func (m *MediaDownloadChanged) ResetMark() {
	m.TitleChanged = false
	m.VisibilityChanged = false
	m.ShareLinkChanged = false

	m.WatchPositionChanged = false
	m.WatchedChanged = false
	m.ViewCountChanged = false
}

func (m *MediaDownloadChanged) MarkUserWatchStatsChanged() {
	m.WatchPositionChanged = true
	m.ViewCountChanged = true
	m.WatchedChanged = true
}

func (m *MediaDownloadChanged) MarkPublicWatchStatsChanged() {
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
