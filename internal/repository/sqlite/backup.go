package sqliterep

import (
	"database/sql"
	"fmt"
)

func (r *Repositories) Backup(dbName string, path string) error {
	backup := func(db *sql.DB, path string) error {
		_, err := db.Exec(`
        VACUUM INTO ?
    `, path)
		return err
	}

	var db *sql.DB = r.dbRegistry.Enttry(dbName).DB()
	if db == nil {
		return fmt.Errorf("unknown db: %s", dbName)
	}

	return backup(db, path)
}
