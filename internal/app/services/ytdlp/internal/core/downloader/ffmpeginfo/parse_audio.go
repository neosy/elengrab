package ffmpeginfo

import (
	"strconv"
	"strings"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

// parseAudio parses an ffmpeg "Audio:" line into AudioInfo.
func (ffi *Info) parseAudio(line string, srcAudioInfo *ddownload.AudioInfo) *ddownload.AudioInfo {
	if line == "" {
		return nil
	}

	var (
		infoSetter audioInfoSetter
		audioInfo  = uptr.Copy(srcAudioInfo)
	)

	if audioInfo == nil {
		audioInfo = &ddownload.AudioInfo{}
	}

	for i, value := range strings.Split(line, ", ") {
		value = strings.TrimSpace(value)

		fields := strings.Fields(value)
		if len(fields) == 0 {
			continue
		}

		// First field is codec
		if i == 0 {
			if len(fields) > 0 {
				audioCodec, err := dtypes.ParseAudioCodec(fields[0])
				if err == nil {
					audioInfo.Codec = audioCodec
					infoSetter.codec = true
				}
			}
			continue
		}

		// Check for sample rate in Hz
		if !infoSetter.sampleRate && strings.Contains(value, " Hz") && ffi.sampleRateRe.MatchString(value) {
			sampleRate, err := strconv.Atoi(fields[0])
			if err == nil && sampleRate != 0 {
				audioInfo.SampleRate = &sampleRate
				infoSetter.sampleRate = true
			}
			continue
		}

		// Check for bitrate in kb/s
		if !infoSetter.bitrate && strings.Contains(value, " kb/s") && ffi.bitrateRe.MatchString(value) {
			bitrate, err := strconv.Atoi(fields[0])
			if err == nil && bitrate != 0 {
				audioInfo.Bitrate = bitrate
				infoSetter.bitrate = true
			}
			continue
		}
	}

	// Discard audio info if codec not set
	if audioInfo != nil && (audioInfo.Codec == "") {
		return nil
	}

	return audioInfo
}
