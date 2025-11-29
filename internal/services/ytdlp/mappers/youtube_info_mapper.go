package mappers

import (
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	dyoutubeinfo "github.com/neosy/elengrab/internal/domain/youtube_info"
	"github.com/neosy/elengrab/internal/services/ytdlp/dto"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

func (m *Mappers) VideoInfoToDomain(info *dto.YouTubeInfo) *dyoutubeinfo.YouTubeInfo {
	var formats = make([]dyoutubeinfo.Format, 0, len(info.Formats))

	for _, f := range info.Formats {
		var fps *int
		if f.FPS > 0 {
			fps = uptr.Int(int(f.FPS))
		}

		var acodec *string
		if f.ACodec != "none" {
			acodec = &f.ACodec
		}

		var vcodec *string
		if f.VCodec != "none" {
			vcodec = &f.VCodec
		}

		var abr *float32
		if f.Abr > 0 {
			abr = &f.Abr
		}

		var vbr *float32
		if f.Vbr > 0 {
			vbr = &f.Vbr
		}

		var formatType = dtypes.FormatTypeNone
		if acodec != nil && vcodec != nil {
			formatType = dtypes.FormatTypeVideoAudio
		} else if acodec == nil {
			formatType = dtypes.FormatTypeVideoOnly
		} else {
			formatType = dtypes.FormatTypeAudioOnly
		}

		formats = append(formats, dyoutubeinfo.Format{
			FormatType: formatType,
			FormatId:   f.FormatID,
			FileExt:    f.Ext,
			Height:     f.Height,
			Width:      f.Width,
			FPS:        fps,
			Format:     f.Format,
			FormatNote: f.FormatNote,
			Resolution: f.Resolution,
			VCodec:     vcodec,
			ACodec:     acodec,
			Vbr:        vbr,
			Abr:        abr,
			Asr:        f.Asr,
			Filesize:   f.Filesize,
		})
	}

	return &dyoutubeinfo.YouTubeInfo{
		Title:   info.Title,
		Formats: formats,
	}
}
