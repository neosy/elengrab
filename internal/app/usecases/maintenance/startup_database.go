package maintenance

import (
	"context"
)

func (m *maintenance) StartupDatabase(ctx context.Context) error {
	return m.repositories.StartupMaintenance(ctx)
}
