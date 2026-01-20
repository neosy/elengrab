package sqliterep

import (
	"database/sql"
	"fmt"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

func (r *Repositories) Backup(dbName persistence.DBName, path string) error {
	backup := func(db *sql.DB, path string) error {
		_, err := db.Exec(`
        VACUUM INTO ?
    `, path)
		return err
	}

	var db *sql.DB = r.dbByName[dbName]
	if db == nil {
		return fmt.Errorf("unknown db: %s", dbName)
	}

	return backup(db, path)
}
