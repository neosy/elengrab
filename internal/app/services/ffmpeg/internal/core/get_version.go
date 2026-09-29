package core

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/neosy/elengrab/internal/app/services/internal/utils"
)

func (c *FFmpegCore) GetFFmpegVersion(ctx context.Context) (string, error) {
	return GetFFmpegVersion(ctx, c.ffmpegPath)
}

func GetFFmpegVersion(ctx context.Context, ffmpegPath string) (string, error) {
	// Prepare command arguments
	var args []string

	args = append(args, "-version")

	// Execute the command to get version
	out, err := utils.ExecCommandContext(ctx, ffmpegPath, args...)
	if err != nil {
		return "", fmt.Errorf("failed to execute %s command: %w", filepath.Base(ffmpegPath), err)
	}

	parts := strings.Fields(string(out))

	if len(parts) < 3 {
		return "", fmt.Errorf("failed to parse ffmpeg version")
	}

	return parts[2], nil
}
