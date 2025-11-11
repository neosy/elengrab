package sqliterep

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"
)

// InitDB opens the SQLite database, applies migrations only if the database file does not exist,
func InitDB(dbPath, migrationsPath string) (*sql.DB, error) {
	// Open the database (SQLite will create the file if it does not exist)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

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

	return db, nil
}
