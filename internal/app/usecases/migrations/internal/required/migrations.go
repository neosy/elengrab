package required

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/app/usecases/migrations/internal/base"
	"github.com/neosy/elengrab/internal/app/usecases/migrations/internal/dependencies"
	"github.com/neosy/elengrab/internal/app/usecases/migrations/internal/registry"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type migrations struct {
	base.Migrations
}

func NewMigrations(
	logger *slog.Logger,
	dlStorage pstorage.DownloadsStorage,
	usecases dependencies.Usecases,
	services dependencies.Services,
) *migrations {

	migrations := &migrations{
		Migrations: base.NewMigrations(logger, dlStorage, usecases, services),
	}

	migrations.initMigrations()

	return migrations
}

func (m *migrations) Add(id string, run registry.MigrationRunner) {
	m.Migrations.Add(id, dtypes.MigrationTypeNoneRequired, run)
}
