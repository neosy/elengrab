package idto

import (
	"sync"

	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

type DownloadMeta struct {
	URL                 string
	Title               string
	Description         string
	FileName            string
	FileExt             string
	FileFullName        string
	FileSize            *int64
	ChannelID           *string
	ChannelURL          string
	ChannelTitle        string
	MediaInfo           *dservices.MediaInfo
	Thumbnail           *dtypes.ImageData
	ThumbnailVideoFrame *dtypes.ImageData
	Channel             *dtypes.Channel
	Progress            *dservices.DownloaderProgress
}

type SafeDownloadMeta struct {
	mu   sync.RWMutex
	Meta *DownloadMeta
}

func (m *SafeDownloadMeta) Lock() {
	m.mu.Lock()
}

func (m *SafeDownloadMeta) Unlock() {
	m.mu.Unlock()
}

func (m *SafeDownloadMeta) RLock() {
	m.mu.RLock()
}

func (m *SafeDownloadMeta) RUnlock() {
	m.mu.RUnlock()
}

func (m *SafeDownloadMeta) CopyMeta() *DownloadMeta {
	m.RLock()
	defer m.RUnlock()

	if m == nil || m.Meta == nil {
		return nil
	}

	metaCopy := *m.Meta
	metaCopy.FileSize = uptr.Copy(m.Meta.FileSize)
	metaCopy.ChannelID = uptr.Copy(m.Meta.ChannelID)
	metaCopy.MediaInfo = m.Meta.MediaInfo.Copy()
	metaCopy.Channel = m.Meta.Channel.Copy()
	metaCopy.Progress = uptr.Copy(m.Meta.Progress)
	metaCopy.Thumbnail = m.Meta.Thumbnail.Copy()
	metaCopy.ThumbnailVideoFrame = m.Meta.ThumbnailVideoFrame.Copy()

	return &metaCopy
}

func (m *SafeDownloadMeta) InitialResult() *dservices.DownloaderResult {
	meta := m.CopyMeta()

	var ext = meta.FileExt
	if ext == "" && meta.MediaInfo != nil && meta.MediaInfo.Format != dtypes.FileFormatNone {
		ext = meta.MediaInfo.Format.String()
	}

	return &dservices.DownloaderResult{
		ChannelID:           meta.ChannelID,
		MediaTitle:          meta.Title,
		MediaDescription:    uptr.NonZeroString(meta.Description),
		Filename:            meta.FileName,
		FileExt:             ext,
		FileFullName:        meta.FileFullName,
		Filesize:            meta.FileSize,
		MediaInfo:           meta.MediaInfo,
		Channel:             meta.Channel,
		Thumbnail:           meta.Thumbnail,
		ThumbnailVideoFrame: meta.ThumbnailVideoFrame,
		Progress:            meta.Progress,
	}
}
