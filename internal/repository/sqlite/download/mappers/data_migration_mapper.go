package mappers

import (
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
)

func (m *Mappers) MapDataMigrationDomainToEntity(migration *ddownload.DataMigration) (*edownload.DataMigration, error) {
	return &edownload.DataMigration{
		MigrationID:   migration.MigrationID,
		MigrationType: migration.MigrationType.String(),
		Description:   migration.Description,
	}, nil
}

func (m *Mappers) MapDataMigrationEntityToDomain(eMigration *edownload.DataMigration) (*ddownload.DataMigration, error) {
	migrationType, err := dtypes.ParseMigrationType(eMigration.MigrationType)
	if err != nil {
		return nil, err
	}

	return &ddownload.DataMigration{
		MigrationID:   eMigration.MigrationID,
		MigrationType: migrationType,
		Description:   eMigration.Description,
		CreatedAt:     eMigration.CreatedAt,
	}, nil
}
