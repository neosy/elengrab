package core

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// ExtractVideoAudioInfoWithFFprobe extracts video and audio information from a media file
// using ffprobe. It returns VideoInfo, AudioInfo, and media duration if available.
//
// Duration is returned in milliseconds.
func (c *FFmpegCore) ExtractVideoAudioInfoWithFFprobe(
	ctx context.Context,
	filePath string,
	srcMediaInfo *dservices.MediaInfo,
) (*dservices.MediaInfo, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path is empty")
	}

	out, err := exec.CommandContext(
		ctx,
		c.ffprobePath,
		"-v", "error",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		filePath,
	).Output()

	if err != nil {
		c.logger.Warn(
			"ffprobe command error",
			"error", err,
		)
		return nil, fmt.Errorf("ffprobe command error: %w", err)
	}

	if len(out) == 0 {
		return nil, fmt.Errorf("ffprobe output is empty")
	}

	// ---- ffprobe response model ----

	type ffprobeFormat struct {
		Duration string `json:"duration"`
	}

	type ffprobeResult struct {
		Format  ffprobeFormat   `json:"format"`
		Streams []ffprobeStream `json:"streams"`
	}

	var probe ffprobeResult
	if err := json.Unmarshal(out, &probe); err != nil {
		return nil, fmt.Errorf("failed to parse ffprobe json: %w", err)
	}

	var (
		duration time.Duration
	)
	if probe.Format.Duration != "" {
		seconds, err := strconv.ParseFloat(probe.Format.Duration, 64)
		if err == nil && seconds > 0 {
			duration = time.Duration(float64(time.Second) * math.Round(seconds*1e6) / 1e6)
		}
	}

	var (
		videoInfo *dtypes.VideoInfo
		audioInfo *dtypes.AudioInfo
	)

	for _, s := range probe.Streams {
		switch s.CodecType {
		case "video":
			if videoInfo == nil {
				videoInfo = c.info.parseVideoFromFFprobe(s, srcMediaInfo.VideoInfo)
			}

		case "audio":
			if audioInfo == nil {
				audioInfo = c.info.parseAudioFromFFprobe(s, srcMediaInfo.AudioInfo)
			}
		}

		if videoInfo != nil && audioInfo != nil {
			break
		}
	}

	var formatType dtypes.FormatType
	if videoInfo != nil && audioInfo != nil {
		formatType = dtypes.FormatTypeVideoAudio
	} else if videoInfo != nil {
		formatType = dtypes.FormatTypeVideoOnly
	} else if audioInfo != nil {
		formatType = dtypes.FormatTypeAudioOnly
	}

	return &dservices.MediaInfo{
		FormatType: formatType,
		Format:     dtypes.MapFileExtToFileFormat(filepath.Ext(filePath)),
		Duration:   duration,
		VideoInfo:  videoInfo,
		AudioInfo:  audioInfo,
	}, nil
}

// ExtractVideoAudioInfoWithFFmpeg extracts video and audio information from a media file
// using ffmpeg. Returns VideoInfo and AudioInfo, or nil if unavailable.
func (c *FFmpegCore) ExtractVideoAudioInfoWithFFmpeg(
	ctx context.Context,
	filePath string,
	srcMediaInfo *dservices.MediaInfo,
) (*dservices.MediaInfo, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path is empty")
	}

	// Run ffmpeg -i to get media info and capture combined stdout/stderr
	// ffmpeg -i returns non-zero exit code, but we can still parse the output
	out, err := exec.CommandContext(
		ctx,
		c.ffmpegPath,
		"-i",
		filePath,
	).CombinedOutput()

	if err != nil && len(out) == 0 {
		c.logger.Warn(
			"ffmpeg command error",
			"error", err,
		)
		return nil, fmt.Errorf("ffmpeg command error: %w", err)
	}

	if len(out) == 0 {
		return nil, fmt.Errorf("ffmpeg output is empty")
	}

	var (
		videoLine, audioLine, durationLine string
	)

	// Extract first lines containing video and audio info
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		if durationLine == "" && strings.Contains(line, "Duration:") {
			_, durationLine, _ = strings.Cut(line, "Duration:")
			durationLine = strings.TrimSpace(durationLine)
		}
		if videoLine == "" && strings.Contains(line, "Video:") {
			_, videoLine, _ = strings.Cut(line, "Video:")
			videoLine = strings.TrimSpace(videoLine)
		}
		if audioLine == "" && strings.Contains(line, "Audio:") {
			_, audioLine, _ = strings.Cut(line, "Audio:")
			audioLine = strings.TrimSpace(audioLine)
		}
		if videoLine != "" && audioLine != "" && durationLine != "" {
			break
		}
	}

	if err != nil && videoLine == "" && audioLine == "" {
		c.logger.Warn(
			"ffmpeg command error",
			"error", err,
			"output", string(out),
		)
		return nil, fmt.Errorf("ffmpeg command error: %w", err)
	}

	// Return nil if no relevant lines found
	if videoLine == "" && audioLine == "" {
		return nil, nil
	}

	// Extract duration
	duration := c.extractDuration(durationLine)

	// Parse video information
	videoInfo := c.info.parseVideoFromFFmppeg(videoLine, srcMediaInfo.VideoInfo)

	// Parse audio information
	audioInfo := c.info.parseAudioFromFFmppeg(audioLine, srcMediaInfo.AudioInfo)

	var formatType dtypes.FormatType
	if videoInfo != nil && audioInfo != nil {
		formatType = dtypes.FormatTypeVideoAudio
	} else if videoInfo != nil {
		formatType = dtypes.FormatTypeVideoOnly
	} else if audioInfo != nil {
		formatType = dtypes.FormatTypeAudioOnly
	}

	return &dservices.MediaInfo{
		FormatType: formatType,
		Format:     dtypes.MapFileExtToFileFormat(filepath.Ext(filePath)),
		Duration:   duration,
		VideoInfo:  videoInfo,
		AudioInfo:  audioInfo,
	}, nil
}

func (c *FFmpegCore) extractDuration(line string) time.Duration {
	if line == "" {
		return 0
	}

	values := strings.Split(line, ", ")
	if len(values) == 0 {
		return 0
	}

	value := values[0]
	if !c.info.durationRe.MatchString(value) {
		return 0
	}

	times := strings.Split(value, ":")
	if len(times) != 3 {
		return 0
	}

	h, err := strconv.Atoi(times[0])
	if err != nil {
		return 0
	}

	m, err := strconv.Atoi(times[1])
	if err != nil {
		return 0
	}

	s, err := strconv.ParseFloat(times[2], 64)
	if err != nil {
		return 0
	}
	s = math.Round(s*1e6) / 1e6

	return time.Duration((float64(h*60*60+m*60) + s) * float64(time.Second))
}
