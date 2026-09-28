package denosrv

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/app/services/deno/internal/consts"
	"github.com/neosy/elengrab/internal/app/services/deno/internal/core"
	"github.com/neosy/elengrab/internal/app/services/deno/internal/utils"
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
	cmdDenoPath, err := utils.ResolveCmdPath(consts.DenoName, binDir)
	if err != nil {
		return nil, err
	}

	err = utils.CheckDeno(consts.DenoName)
	if err != nil {
		return nil, err
	} else {
		logger.Info("Deno executable found in PATH", "executable", consts.DenoName)
	}

	return &DenoService{
		logger: logger,
		core:   core.NewDenoCore(logger, cmdDenoPath),
	}, nil
}
