package ffmpeg

import (
	"bufio"
	"bytes"
	"context"
	"os/exec"
	"strings"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/utils"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

// GetVideoAudioInfoFromFile extracts video and audio information from a media file
// using ffmpeg. Returns VideoInfo and AudioInfo, or nil if unavailable.
func (ffi *Info) GetVideoAudioInfoFromFile(
	ctx context.Context,
	filePath string,
	srcMediaInfo *ddownload.MediaInfo,
) (*ddownload.VideoInfo, *ddownload.AudioInfo) {
	if filePath == "" {
		return nil, nil
	}

	// Resolve ffmpeg executable path
	ffmpegPath, err := utils.ResolveCmdPath(ffmpegName, "")
	if err != nil {
		return nil, nil
	}

	// Run ffmpeg -i to get media info and capture combined stdout/stderr
	out, _ := exec.CommandContext(
		ctx,
		ffmpegPath,
		"-i",
		filePath,
	).CombinedOutput()

	if len(out) == 0 {
		return nil, nil
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

	// Return nil if no relevant lines found
	if videoLine == "" && audioLine == "" {
		return nil, nil
	}

	// Parse video information
	videoInfo := ffi.parseVideo(videoLine, srcMediaInfo.VideoInfo)

	// Parse audio information
	audioInfo := ffi.parseAudio(audioLine, srcMediaInfo.AudioInfo)

	return videoInfo, audioInfo
}
