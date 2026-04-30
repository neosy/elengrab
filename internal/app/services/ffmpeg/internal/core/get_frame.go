package core

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/imgx"
)

func (c *FFmpegCore) GetFrame(
	ctx context.Context,
	filePath string,
	args ...string,
) (*dtypes.ImageData, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path is empty")
	}

	newArgs := append(
		[]string{
			"-i", filePath,
		},
		args...,
	)
	newArgs = append(newArgs, "-")

	cmd := exec.CommandContext(ctx, c.ffmpegPath, newArgs...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil && stdout.Len() == 0 {
		c.logger.Warn(
			"ffmpeg command error",
			"error", err,
			"stderr", stderr.String(),
			"args", newArgs,
		)
		return nil, fmt.Errorf("ffmpeg command error: %w", err)
	}

	raw := stdout.Bytes()

	if len(raw) == 0 {
		return nil, fmt.Errorf("data is empty")
	}

	size, err := imgx.ImageSize(raw)
	if err != nil {
		c.logger.Warn(
			"Failed to decode image data",
			"error", err,
		)
	}

	return &dtypes.ImageData{
		Format: dtypes.ImageFormatJPEG,
		Width:  size.Width,
		Height: size.Height,
		Raw:    raw,
	}, nil
}
