package pworkers

import "context"

type AuthWebMaintenanceRunner interface {
	Startup(ctx context.Context) error
}
