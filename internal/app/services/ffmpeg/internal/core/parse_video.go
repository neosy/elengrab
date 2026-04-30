package core

import (
	"strconv"
	"strings"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

// parseVideo parses an ffmpeg "Video:" line into VideoInfo.
func (info *info) parseVideo(line string, srcVideoInfo *dtypes.VideoInfo) *dtypes.VideoInfo {
	if line == "" {
		return nil
	}

	var (
		infoSetter videoInfoSetter
		videoInfo  = uptr.Copy(srcVideoInfo)
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
		if i == 0 {
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
	if videoInfo != nil && (videoInfo.Codec == "") {
		return nil
	}

	return videoInfo
}
