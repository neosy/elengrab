package executor

import (
	"context"
	"fmt"
	"strings"

	"github.com/neosy/elengrab/internal/app/services/internal/utils"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
)

func (e *Executor) GetVersion(ctx context.Context) (string, error) {
	return GetVersion(ctx, e.ytDlpPath)
}

func GetVersion(ctx context.Context, ytDlpPath string) (string, error) {
	// Prepare command arguments
	var args []string

	args = append(args, "--version")

	// Execute the command to get version
	out, err := utils.ExecCommandContext(ctx, ytDlpPath, args...)
	if err != nil {
		return "", fmt.Errorf("failed to execute %s command: %w", consts.YtDlpName, err)
	}

	return strings.TrimSpace(string(out)), nil
}
