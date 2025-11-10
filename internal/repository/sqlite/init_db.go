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

// InitDB opens the SQLite database, applies migrations only if the database file does not exist,
func InitDB(dbPath, migrationsPath string) (*sql.DB, error) {
	// Check if the database file exists
	_, err := os.Stat(dbPath)
	isNewDB := os.IsNotExist(err)

	// Open the database (SQLite will create the file if it does not exist)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Apply migrations only if this is a new database
	if isNewDB {
		// Create a driver instance for migrate
		driver, err := sqlite.WithInstance(db, &sqlite.Config{})
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to create migrate driver: %w", err)
		}

		// Create a migrate instance pointing to the migrations folder
		m, err := migrate.NewWithDatabaseInstance(
			migrationsPath,
			"sqlite",
			driver,
		)
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to create migrate instance: %w", err)
		}

		// Apply all up migrations
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			db.Close()
			return nil, fmt.Errorf("migration failed: %w", err)
		}

		log.Println("Database created and migrations applied successfully")
	}

	return db, nil
}
