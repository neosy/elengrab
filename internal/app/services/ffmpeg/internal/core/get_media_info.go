package core

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// GetVideoAudioInfo extracts video and audio information from a media file
// using ffmpeg. Returns VideoInfo and AudioInfo, or nil if unavailable.
func (c *FFmpegCore) GetVideoAudioInfo(
	ctx context.Context,
	filePath string,
	srcMediaInfo *dservices.MediaInfo,
) (*dtypes.VideoInfo, *dtypes.AudioInfo, error) {
	if filePath == "" {
		return nil, nil, fmt.Errorf("file path is empty")
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
		return nil, nil, fmt.Errorf("ffmpeg command error: %w", err)
	}

	if len(out) == 0 {
		return nil, nil, fmt.Errorf("ffmpeg output is empty")
	}

	var (
		videoLine, audioLine string
	)

	// Extract first lines containing video and audio info
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		if videoLine == "" && strings.Contains(line, "Video:") {
			_, videoLine, _ = strings.Cut(line, "Video:")
			videoLine = strings.TrimSpace(videoLine)
		}
		if audioLine == "" && strings.Contains(line, "Audio:") {
			_, audioLine, _ = strings.Cut(line, "Audio:")
			audioLine = strings.TrimSpace(audioLine)
		}
		if videoLine != "" && audioLine != "" {
			break
		}
	}

	if err != nil && videoLine == "" && audioLine == "" {
		c.logger.Warn(
			"ffmpeg command error",
			"error", err,
			"output", string(out),
		)
		return nil, nil, fmt.Errorf("ffmpeg command error: %w", err)
	}

	// Return nil if no relevant lines found
	if videoLine == "" && audioLine == "" {
		return nil, nil, nil
	}

	// Parse video information
	videoInfo := c.info.parseVideo(videoLine, srcMediaInfo.VideoInfo)

	// Parse audio information
	audioInfo := c.info.parseAudio(audioLine, srcMediaInfo.AudioInfo)

	return videoInfo, audioInfo, nil
}
