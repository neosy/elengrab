package mappers

import (
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
)

func (m *Mappers) MapDataMigrationDomainToEntity(migration *ddownload.DataMigration) (*edownload.DataMigration, error) {
	return &edownload.DataMigration{
		MigrationID: migration.MigrationID,
		Description: migration.Description,
	}, nil
}

func (m *Mappers) MapDataMigrationEntityToDomain(eMigration *edownload.DataMigration) (*ddownload.DataMigration, error) {
	return &ddownload.DataMigration{
		MigrationID: eMigration.MigrationID,
		Description: eMigration.Description,
		CreatedAt:   eMigration.CreatedAt,
	}, nil
}
