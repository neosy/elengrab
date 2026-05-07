package mappers

import (
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

func (m *Mappers) MapMediaInfoToDomain(info *idto.MediaInfo) *dmedia.MediaInfo {
	var formats = make([]dmedia.MediaFormat, 0, len(info.Formats))

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

		formats = append(formats, dmedia.MediaFormat{
			FormatType: formatType,
			FormatId:   f.FormatID,
			FileExt:    f.FileExt,
			Height:     int(f.Height),
			Width:      int(f.Width),
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

	return &dmedia.MediaInfo{
		Title:       info.Title,
		Description: info.Description,
		Formats:     formats,
	}
}
