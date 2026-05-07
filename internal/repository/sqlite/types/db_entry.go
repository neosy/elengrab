package sqlitetypes

import (
	"database/sql"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type dbEntry struct {
	schema persistence.DBSchema
	db     *sql.DB
}

func NewDBEntry(schema persistence.DBSchema, db *sql.DB) persistence.DBEntry {
	return &dbEntry{
		schema: schema,
		db:     db,
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
