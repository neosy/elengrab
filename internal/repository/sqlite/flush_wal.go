package sqliterep

import "database/sql"

// flushWAL flushes all WAL changes into the main SQLite database file.
func flushWAL(db *sql.DB) error {
	_, err := db.Exec("PRAGMA wal_checkpoint(FULL);")
	return err
}

func (r *Repositories) FlushWAL() error {
	for _, dbName := range r.GetDBNames() {
		err := flushWAL(r.dbByName[dbName])
		if err != nil {
			return err
		}
	}

	return nil
}
