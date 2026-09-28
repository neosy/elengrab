package executor

import (
	"context"
	"fmt"
	"strings"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/utils"
)

func (e *Executor) GetVersion(ctx context.Context) (string, error) {
	// Prepare command arguments
	var args []string

	args = append(args, "--version")

	// Execute the command to get version
	out, err := utils.ExecCommandContext(ctx, e.ytDlpPath, args...)
	if err != nil {
		return "", fmt.Errorf("failed to execute %s command: %w", consts.YtDlpName, err)
	}

	return strings.TrimSpace(string(out)), nil
}
