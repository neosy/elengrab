package sqlitetypes

import (
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type DBRegistry struct {
	schemas       []persistence.DBSchema
	schemasByName map[string]persistence.DBSchema
	entriesByName map[string]persistence.DBEntry
}

func NewRegistry(entries map[string]persistence.DBEntry) *DBRegistry {
	reg := &DBRegistry{
		schemas:       make([]persistence.DBSchema, 0, len(entries)),
		schemasByName: make(map[string]persistence.DBSchema, len(entries)),
		entriesByName: make(map[string]persistence.DBEntry, len(entries)),
	}

	for _, e := range entries {
		reg.add(e)
	}

	return reg
}

func (r *DBRegistry) add(dbEntry persistence.DBEntry) {
	r.schemas = append(r.schemas, dbEntry.Schema())
	r.schemasByName[dbEntry.DBName()] = dbEntry.Schema()
	r.entriesByName[dbEntry.DBName()] = dbEntry
}

func (r *DBRegistry) Schemas() []persistence.DBSchema {
	return r.schemas
}

func (r *DBRegistry) SchemasByName() map[string]persistence.DBSchema {
	return r.schemasByName
}

func (r *DBRegistry) EntriesByName() map[string]persistence.DBEntry {
	return r.entriesByName
}

func (r *DBRegistry) Enttry(name string) persistence.DBEntry {
	return r.entriesByName[name]
}
