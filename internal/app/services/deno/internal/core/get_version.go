package core

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/neosy/elengrab/internal/app/services/deno/internal/consts"
	"github.com/neosy/elengrab/internal/app/services/internal/utils"
)

func (c *DenoCore) GetVersion(ctx context.Context) (string, error) {
	return GetVersion(ctx, c.denoPath)
}

func GetVersion(ctx context.Context, denoPath string) (string, error) {
	// Prepare command arguments
	var args []string

	args = append(args, "--version")

	// Execute the command to get version
	out, err := utils.ExecCommandContext(ctx, denoPath, args...)
	if err != nil {
		return "", fmt.Errorf("failed to execute %s command: %w", consts.DenoName, err)
	}

	parts := strings.Fields(string(out))

	if len(parts) < 2 {
		return "", errors.New("invalid Deno version output")
	}

	return parts[1], nil
}
