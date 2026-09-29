package utils

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
)

func ExecCommandContext(ctx context.Context, name string, arg ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, arg...)

	// Buffers to capture stdout and stderr
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// The context was canceled
		if ctx.Err() != nil {
			return nil, fmt.Errorf("process canceled: %w", ctx.Err())
		}
		errOut := fmt.Errorf("%s failed: %v, stderr: %s", filepath.Base(name), err, stderr.String())
		return nil, errOut
	}

	return out.Bytes(), nil
}
