package denosrv

import (
	"log/slog"
	"path/filepath"

	"github.com/neosy/elengrab/internal/app/services/deno/internal/consts"
	"github.com/neosy/elengrab/internal/app/services/deno/internal/core"
	"github.com/neosy/elengrab/internal/app/services/internal/utils"
)

// DenoService represents a service for interacting with deno.
type DenoService struct {
	logger *slog.Logger

	// internal
	core *core.DenoCore
}

// NewDenoService creates a new instance of DenoService.
func NewDenoService(
	logger *slog.Logger,
	binDir string,
) (*DenoService, error) {
	denoPath, err := utils.ResolveCmdPath(consts.DenoName, binDir)
	if err != nil {
		return nil, err
	}

	err = checkDeno(denoPath)
	if err != nil {
		return nil, err
	}

	logger.Info("Deno executable found",
		"name", filepath.Base(denoPath),
		"dir", filepath.Dir(denoPath),
	)

	return &DenoService{
		logger: logger,
		core:   core.NewDenoCore(logger, denoPath),
	}, nil
}
