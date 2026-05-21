package core

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// ExtractDurationMsWithFFprobe extracts media duration from the given file using ffprobe.
//
// It returns duration in milliseconds. If ffprobe fails or duration cannot be parsed,
// an error is returned.
func (c *FFmpegCore) ExtractDurationMsWithFFprobe(ctx context.Context, filePath string) (time.Duration, error) {
	out, err := exec.CommandContext(
		ctx,
		c.ffprobePath,
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filePath,
	).Output()

	if err != nil {
		return 0, fmt.Errorf("ffprobe failed: %w", err)
	}

	s := strings.TrimSpace(string(out))
	if s == "" {
		return 0, fmt.Errorf("empty duration from ffprobe")
	}

	seconds, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid duration format: %w", err)
	}

	return time.Duration(seconds * float64(time.Second)), nil
}
