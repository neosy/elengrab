package denosrv

import (
	"context"
)

func (srv *DenoService) GetVersion(ctx context.Context) (string, error) {
	return srv.core.GetVersion(ctx)
}
