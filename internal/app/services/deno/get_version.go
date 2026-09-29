package denosrv

import (
	"context"

	"github.com/neosy/elengrab/internal/app/services/deno/internal/consts"
	"github.com/neosy/elengrab/internal/app/services/deno/internal/core"
	"github.com/neosy/elengrab/internal/app/services/internal/utils"
)

func (srv *DenoService) GetVersion(ctx context.Context) (string, error) {
	return srv.core.GetVersion(ctx)
}

// GetVersion retrieves the version of the Deno executable.
func GetVersion(ctx context.Context) (string, error) {
	cmdPath, err := utils.ResolveCmdPath(consts.DenoName, "")
	if err != nil {
		return "", err
	}

	return core.GetVersion(ctx, cmdPath)
}
