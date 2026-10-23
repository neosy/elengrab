package downloadpreparer

import (
	"context"
	"fmt"
	"strings"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

func (p *DownloadPreparer) Prepare(
	ctx context.Context,
	url string,
	dlOptions idto.DLOptions,
) (DownloadPlan, error) {
	var extractResult extractResult

	switch dlOptions.FormatType {
	// Video + Audio or Video only
	case dtypes.FormatTypeVideoAudio, dtypes.FormatTypeVideoOnly:
		var err error
		extractResult, err = p.extractVideoAudio(ctx, url, dlOptions)
		if err != nil {
			return DownloadPlan{}, err
		}
	// Audio only
	case dtypes.FormatTypeAudioOnly:
		var err error
		extractResult, err = p.extractAudioOnly(ctx, url, dlOptions)
		if err != nil {
			return DownloadPlan{}, err
		}
	default:
		return DownloadPlan{}, fmt.Errorf("unsupported format type: %q", dlOptions.FormatType)
	}

	args := make([]string, 0)

	// prevent downloading the entire playlist, only fetch single video
	args = append(args, "--no-playlist")
	args = append(args, "--no-warnings")

	// Add format option to yt-dlp arguments
	if len(extractResult.selectedFormatIDs) > 0 {
		args = append(args, "-f", strings.Join(extractResult.selectedFormatIDs, "+"))
	} else {
		args = append(args, "-f", extractResult.formatQuery)
	}

	// Add yt-dlp arguments
	args = append(args, extractResult.args...)

	downloadResult := DownloadPlan{
		Args:         args,
		FileExtNoDot: extractResult.fileExtNoDot,
		ExtractInfo:  extractResult.info,
		MediaInfo: &pservices.MediaInfo{
			FormatType: dlOptions.FormatType,
			Format:     dtypes.FileExtToFileFormat(extractResult.fileExtNoDot),
			VideoInfo:  nil,
			AudioInfo:  nil,
		},
	}

	if extractResult.info != nil && len(extractResult.info.Formats) > 0 && extractResult.mediaFormat != nil {
		if extractResult.videoCodec != dtypes.VideoCodecNone {
			downloadResult.MediaInfo.VideoInfo = &dtypes.VideoInfo{
				Codec: extractResult.videoCodec,
				Resolution: dtypes.ParseVideoResolutionWH(
					uint16(extractResult.mediaFormat.Width),
					uint16(extractResult.mediaFormat.Height),
				),
				Bitrate: int(extractResult.mediaFormat.Vbr),
				Width:   int(extractResult.mediaFormat.Width),
				Height:  int(extractResult.mediaFormat.Height),
			}
		}
		if extractResult.audioCodec != dtypes.AudioCodecNone {
			downloadResult.MediaInfo.AudioInfo = &dtypes.AudioInfo{
				Codec:      extractResult.audioCodec,
				Bitrate:    int(extractResult.mediaFormat.Abr),
				SampleRate: extractResult.mediaFormat.Asr,
			}
		}
	}

	return downloadResult, nil
}
