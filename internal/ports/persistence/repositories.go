package persistence

import (
	"context"
	"database/sql"
)

type DBEntry interface {
	Schema() DBSchema
	DB() *sql.DB
	DBName() string
}

type DBSchema interface {
	DBName() string
	Path(dir string) string
	FileName() string
	TableNames() []string
}

// Repositories is an interface that defines methods for interacting with databases.
type Repositories interface {
	Backup(dbName string, path string) error
	FlushWAL() error
	Schemas() []DBSchema
	StartupMaintenance(ctx context.Context) error
}
