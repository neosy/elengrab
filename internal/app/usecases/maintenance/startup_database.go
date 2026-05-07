package maintenance

import (
	"context"
)

func (m *Maintenance) StartupDatabase(ctx context.Context) error {
	return m.repositories.StartupMaintenance(ctx)
}
