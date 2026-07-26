package sqliterep

import "database/sql"

// flushWAL flushes all WAL changes into the main SQLite database file.
func flushWAL(db *sql.DB) error {
	_, err := db.Exec("PRAGMA wal_checkpoint(FULL);")
	return err
}

// flushAndTruncateWAL flushes all WAL changes and truncates the WAL file.
func flushAndTruncateWAL(db *sql.DB) error {
	_, err := db.Exec("PRAGMA wal_checkpoint(TRUNCATE);")
	return err
}

func (r *Repositories) FlushWAL() error {
	for _, entry := range r.EntriesByName() {
		err := flushWAL(entry.DB())
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repositories) FlushAndTruncateWAL() error {
	for _, entry := range r.EntriesByName() {
		err := flushAndTruncateWAL(entry.DB())
		if err != nil {
			return err
		}
	}

	return nil
}
