package sqliterep

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"
)

const (
	// Journal mode for SQLite: Write-Ahead Logging enables concurrent reads during writes
	sqliteJournalMode = "WAL"
	// Busy timeout in milliseconds: how long SQLite waits for a lock before returning an error
	sqliteBusyTimeoutMS = 5000
	// Synchronous mode: balances write speed and reliability (NORMAL is usually sufficient with WAL)
	sqliteSynchronous = "NORMAL"
)

// InitDB opens the SQLite database with WAL mode and applies migrations.
// dbPath - path to SQLite file
// migrationsPath - path to folder with migrations (e.g., "./db/migrations")
func InitDB(dbPath, migrationsDir string) (*sql.DB, error) {
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

	// 4. Create SQLite driver instance for golang-migrate
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create migrate driver: %w", err)
	}

	// 5. Check if migrations directory exists
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		db.Close()
		return nil, fmt.Errorf("migrations directory does not exist: %s", migrationsDir)
	}

	// 6. Create migrate instance pointing to the migrations folder
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsDir),
		"sqlite",
		driver,
	)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// ******** rollback migration
	// if err := m.Force(2); err != nil {
	// 	log.Fatalf("failed to force version: %v", err)
	// }
	// if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
	// 	log.Fatalf("failed to rollback migration: %v", err)
	// }

	// 7. Apply all up migrations
	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Println("All migrations are already applied. Database is up to date.")
		} else {
			db.Close()
			log.Printf("migration failed: %v\n", err)
			return nil, fmt.Errorf("migration failed: %w", err)
		}
	} else {
		log.Println("New migrations applied successfully.")
	}

	// 8. Return the live database connection
	return db, nil
}
