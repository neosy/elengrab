package core

import (
	"strconv"
	"strings"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// parseAudioFromFFprobe parses an FFprobe "stream" line into AudioInfo.
func (info *info) parseAudioFromFFprobe(stream ffprobeStream, srcAudioInfo *dtypes.AudioInfo) *dtypes.AudioInfo {
	var (
		audioInfo = srcAudioInfo.Copy()
	)

	if audioInfo == nil {
		audioInfo = &dtypes.AudioInfo{}
	}

	audioCodec, err := dtypes.ParseAudioCodec(stream.CodecName)
	if err == nil {
		audioInfo.Codec = audioCodec
	}

	sampleRate, err := strconv.Atoi(stream.SampleRate)
	if err == nil && sampleRate != 0 {
		audioInfo.SampleRate = &sampleRate
	}

	bitrate, err := strconv.Atoi(stream.BitRate)
	if err == nil && bitrate != 0 {
		audioInfo.Bitrate = int(bitrate / 1000)
	}

	// Discard audio info if codec not set
	if audioInfo.Codec == "" {
		return nil
	}

	return audioInfo
}

// parseAudioFromFFmppeg parses an FFmpeg "Audio:" line into AudioInfo.
func (info *info) parseAudioFromFFmppeg(line string, srcAudioInfo *dtypes.AudioInfo) *dtypes.AudioInfo {
	if line == "" {
		return nil
	}

	var (
		infoSetter audioInfoSetter
		audioInfo  = srcAudioInfo.Copy()
	)

	if audioInfo == nil {
		audioInfo = &dtypes.AudioInfo{}
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
		if !infoSetter.sampleRate && strings.Contains(value, " Hz") && info.sampleRateRe.MatchString(value) {
			sampleRate, err := strconv.Atoi(fields[0])
			if err == nil && sampleRate != 0 {
				audioInfo.SampleRate = &sampleRate
				infoSetter.sampleRate = true
			}
			continue
		}

		// Check for bitrate in kb/s
		if !infoSetter.bitrate && strings.Contains(value, " kb/s") && info.bitrateRe.MatchString(value) {
			bitrate, err := strconv.Atoi(fields[0])
			if err == nil && bitrate != 0 {
				audioInfo.Bitrate = bitrate
				infoSetter.bitrate = true
			}
			continue
		}
	}

	// Discard audio info if codec not set
	if audioInfo.Codec == "" {
		return nil
	}

	return audioInfo
}
