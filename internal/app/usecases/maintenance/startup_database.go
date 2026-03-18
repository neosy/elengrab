package maintenance

import (
	"context"
)

func (m *Maintenance) StartupDatabase(ctx context.Context) error {
	return m.database.StartupMaintenance(ctx)
}
