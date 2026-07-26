package sqliterep

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const (
	// Journal mode for SQLite: Write-Ahead Logging enables concurrent reads during writes
	sqliteJournalMode = "WAL"
	// Busy timeout in milliseconds: how long SQLite waits for a lock before returning an error
	sqliteBusyTimeoutMS = 5000
	// Synchronous mode: balances write speed and reliability (NORMAL is usually sufficient with WAL)
	sqliteSynchronous = "NORMAL"
	// Enable SQLite foreign key constraints (required for ON DELETE CASCADE)
	sqliteForeignKeys = "ON"
)

// InitDB opens the SQLite database with WAL mode
// dbPath - path to SQLite file
func InitDB(logger *slog.Logger, dbPath string) (*sql.DB, error) {
	// 1. Build DSN with SQLite parameters:
	// - WAL for concurrent reads
	// - busy_timeout to wait for locks
	// - synchronous=NORMAL for faster writes while maintaining reasonable durability
	dsn := fmt.Sprintf(
		"file:%s?_journal_mode=%s&_busy_timeout=%d&_synchronous=%s",
		dbPath,
		sqliteJournalMode,
		sqliteBusyTimeoutMS,
		sqliteSynchronous,
	)

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// SQLite works best with a single connection.
	// Limiting the pool prevents lock contention and SQLITE_BUSY errors.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// Keep the connection alive indefinitely to avoid reapplying PRAGMAs.
	db.SetConnMaxLifetime(0)

	// 2. Test database connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// 3. Explicitly apply PRAGMA settings to ensure the desired configuration
	if _, err := db.Exec(fmt.Sprintf("PRAGMA journal_mode = %s;", sqliteJournalMode)); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable journal mode: %w", err)
	}
	if _, err := db.Exec(fmt.Sprintf("PRAGMA busy_timeout = %d;", sqliteBusyTimeoutMS)); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set busy_timeout: %w", err)
	}
	if _, err := db.Exec(fmt.Sprintf("PRAGMA synchronous = %s;", sqliteSynchronous)); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set synchronous mode: %w", err)
	}
	if _, err := db.Exec(fmt.Sprintf("PRAGMA foreign_keys = %s;", sqliteForeignKeys)); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set foreign_keys: %w", err)
	}

	// Test WAL
	row := db.QueryRow("PRAGMA journal_mode;")
	var mode string
	if err := row.Scan(&mode); err != nil {
		logger.Warn("failed to read journal_mode:", "error", err)
	}
	logger.Debug(
		fmt.Sprintf("SQLite initialized: PRAGMA journal_mode=%s", mode),
		"Database file", filepath.Base(dbPath),
	)

	// Verify foreign_keys
	row = db.QueryRow("PRAGMA foreign_keys;")
	if err := row.Scan(&mode); err != nil {
		logger.Warn("failed to read foreign_keys:", "error", err)
	}
	if mode != "1" {
		logger.Warn(fmt.Sprintf("failed to set PRAGMA foreign_keys = %s", sqliteForeignKeys))
		db.Close()
		return nil, fmt.Errorf("failed to set foreign_keys = %s", sqliteForeignKeys)
	}

	// 4. Return the live database connection
	return db, nil
}

func CloseDB(db *sql.DB) {
	// 1. Do a full checkpoint to flush WAL into main DB
	if err := flushAndTruncateWAL(db); err != nil {
		log.Printf("Failed to checkpoint WAL: %v", err)
	}

	// 2. Close the database connection
	if err := db.Close(); err != nil {
		log.Printf("Failed to close SQLite database: %v", err)
	}
}

func (r *Repositories) CloseAllDB() {
	for _, entry := range r.EntriesByName() {
		CloseDB(entry.DB())
	}
}
