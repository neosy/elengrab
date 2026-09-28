package core

import (
	"context"
	"fmt"
	"strings"

	"github.com/neosy/elengrab/internal/app/services/ffmpeg/internal/consts"
	"github.com/neosy/elengrab/internal/app/services/ffmpeg/internal/utils"
)

func (c *FFmpegCore) GetFFmpegVersion(ctx context.Context) (string, error) {
	// Prepare command arguments
	var args []string

	args = append(args, "-version")

	// Execute the command to get version
	out, err := utils.ExecCommandContext(ctx, c.ffmpegPath, args...)
	if err != nil {
		return "", fmt.Errorf("failed to execute %s command: %w", consts.FFmpegName, err)
	}

	parts := strings.Fields(string(out))

	if len(parts) < 3 {
		return "", fmt.Errorf("failed to parse ffmpeg version")
	}

	return parts[2], nil
}
