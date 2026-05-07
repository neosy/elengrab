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
	Progress            *dservices.DownloadProgress
	Options             DownloadOptions
}

type DownloadOptions struct {
	ConcurrentFragments    uint8
	RequiresYouTubeCookies bool
	Extractor              string
	ExtractorArgs          *string
	Args                   []string
}

type SafeDownloadMeta struct {
	mu   sync.Mutex
	Meta *DownloadMeta
}

func (m *SafeDownloadMeta) Lock() {
	m.mu.Lock()
}

func (m *SafeDownloadMeta) Unlock() {
	m.mu.Unlock()
}

func (m *SafeDownloadMeta) CopyMeta() *DownloadMeta {
	m.Lock()
	defer m.Unlock()

	if m == nil || m.Meta == nil {
		return nil
	}

	metaCopy := *m.Meta
	metaCopy.FileSize = uptr.Copy(m.Meta.FileSize)
	metaCopy.ChannelID = uptr.Copy(m.Meta.ChannelID)
	metaCopy.MediaInfo = m.Meta.MediaInfo.Copy()
	metaCopy.Channel = m.Meta.Channel.Copy()
	metaCopy.Progress = uptr.Copy(m.Meta.Progress)
	metaCopy.Options.ExtractorArgs = uptr.Copy(m.Meta.Options.ExtractorArgs)
	metaCopy.Options.Args = m.Meta.Options.Args[0:len(m.Meta.Options.Args):len(m.Meta.Options.Args)]
	metaCopy.Thumbnail = m.Meta.Thumbnail.Copy()
	metaCopy.ThumbnailVideoFrame = m.Meta.ThumbnailVideoFrame.Copy()

	return &metaCopy
}

func (m *SafeDownloadMeta) InitialResult() *dservices.DownloadResult {
	meta := m.CopyMeta()

	var ext = meta.FileExt
	if ext == "" && meta.MediaInfo != nil && meta.MediaInfo.Format != dtypes.FileFormatNone {
		ext = meta.MediaInfo.Format.String()
	}

	return &dservices.DownloadResult{
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
