package core

import (
	"strconv"
	"strings"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// parseVideoFromFFprobe parses an FFprobe "stream" line into VideoInfo.
func (info *info) parseVideoFromFFprobe(stream ffprobeStream, srcVideoInfo *dtypes.VideoInfo) *dtypes.VideoInfo {
	var (
		videoInfo = srcVideoInfo.Copy()
	)

	if videoInfo == nil {
		videoInfo = &dtypes.VideoInfo{}
	}

	videoCodec, err := dtypes.ParseVideoCodec(stream.CodecName)
	if err == nil {
		videoInfo.Codec = videoCodec
	}

	if stream.Width > 0 {
		videoInfo.Width = stream.Width
	}

	if stream.Height > 0 {
		videoInfo.Height = stream.Height
	}

	videoInfo.Resolution = dtypes.ParseVideoResolutionWH(uint16(videoInfo.Width), uint16(videoInfo.Height))

	bitrate, err := strconv.Atoi(stream.BitRate)
	if err == nil && bitrate != 0 {
		videoInfo.Bitrate = int(bitrate / 1000)
	}

	// Discard video info if codec not set
	if videoInfo.Codec == "" {
		return nil
	}

	return videoInfo
}

// parseVideo parses an FFmpeg "Video:" line into VideoInfo.
func (info *info) parseVideoFromFFmppeg(line string, srcVideoInfo *dtypes.VideoInfo) *dtypes.VideoInfo {
	if line == "" {
		return nil
	}

	var (
		infoSetter videoInfoSetter
		videoInfo  = srcVideoInfo.Copy()
	)

	if videoInfo == nil {
		videoInfo = &dtypes.VideoInfo{}
	}

	for i, value := range strings.Split(line, ", ") {
		value = strings.TrimSpace(value)

		fields := strings.Fields(value)
		if len(fields) == 0 {
			continue
		}

		// First field is codec
		if !infoSetter.codec && i == 0 {
			videoCodec, err := dtypes.ParseVideoCodec(fields[0])
			if err == nil {
				videoInfo.Codec = videoCodec
				infoSetter.codec = true
			}
			continue
		}

		// Check for resolution in format WIDTHxHEIGHT
		if !infoSetter.resolution && info.resolutionRe.MatchString(fields[0]) {
			res := fields[0]
			w, h, err := dtypes.VideoResolutionStringToWH(res)
			if err == nil && w != 0 && h != 0 {
				videoInfo.Width = int(w)
				videoInfo.Height = int(h)
				infoSetter.resolution = true
			}
			resolution, err := dtypes.ParseVideoResolutionFromString(res)
			if err == nil {
				videoInfo.Resolution = resolution
				infoSetter.resolution = true
			}
			continue
		}

		// Check for bitrate in kb/s
		if !infoSetter.bitrate && strings.Contains(value, " kb/s") && info.bitrateRe.MatchString(value) {
			bitrate, err := strconv.Atoi(fields[0])
			if err == nil && bitrate != 0 {
				videoInfo.Bitrate = bitrate
				infoSetter.bitrate = true
			}
			continue
		}
	}

	// Discard video info if codec not set
	if videoInfo.Codec == "" {
		return nil
	}

	return videoInfo
}
