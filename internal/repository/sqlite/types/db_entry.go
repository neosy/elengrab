package types

import (
	"database/sql"

	"github.com/neosy/elengrab/internal/ports/persistence"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
)

type dbEntry struct {
	schema persistence.DBSchema
	db     *sql.DB
	locker dbexec.WriteLocker
}

func NewDBEntry(schema persistence.DBSchema, db *sql.DB) persistence.DBEntry {
	return &dbEntry{
		schema: schema,
		db:     db,
		locker: dbexec.NewSQLiteLock(),
	}
}

func (e *dbEntry) Schema() persistence.DBSchema {
	return e.schema
}

func (e *dbEntry) DB() *sql.DB {
	return e.db
}

func (e *dbEntry) DBName() string {
	return e.schema.DBName()
}

func (e *dbEntry) Locker() dbexec.WriteLocker {
	return e.locker
}
