package sqliterep

import "database/sql"

// FlushWAL flushes all WAL changes into the main SQLite database file.
func FlushWAL(db *sql.DB) error {
	_, err := db.Exec("PRAGMA wal_checkpoint(FULL);")
	return err
}

func (r *Repositories) FlushWAL() error {
	return FlushWAL(r.db)
}
