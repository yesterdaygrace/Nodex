package models

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the shared database handle used across the app.
var DB *gorm.DB

// defaultDSN is the PostgreSQL connection string used when no
// DATABASE_URL environment variable is set. It matches the local
// practice setup: a Postgres instance exposing port 5435 with the
// "postgres" user/password and a "test_db" database.
const defaultDSN = "host=127.0.0.1 user=postgres password=postgres dbname=test_db port=5435 sslmode=disable"

// dbDSN resolves the connection string from the environment, falling
// back to the local default above.
//
// This follows the 12-factor guideline "store config in the environment":
// the database location is not baked into the binary, so you can point the
// app at a differently-ported or remote Postgres without recompiling.
func dbDSN() string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}
	return defaultDSN
}

// ConnectDatabase opens the database, runs the schema migrations, and
// assigns the resulting handle to DB.
//
// It returns the handle and any error instead of panicking, so the caller
// (main.go) can decide how to handle startup failure. The global DB is
// also set for backward compat during the migration — callers still reading
// models.DB keep working while new code can use the returned *gorm.DB.
//
// Lifecycle events are logged via log/slog as structured records.
func ConnectDatabase() (*gorm.DB, error) {
	dsn := dbDSN()

	const attempts = 5
	var lastErr error

	for i := 1; i <= attempts; i++ {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			// Connection succeeded: migrate the schema, then publish the handle.
			if err := db.AutoMigrate(&Post{}); err != nil {
				return nil, fmt.Errorf("failed to migrate database: %w", err)
			}
			DB = db
			slog.Info("database connection established")
			return db, nil
		}

		lastErr = err
		slog.Debug("database not ready",
			slog.Int("attempt", i),
			slog.Int("of", attempts),
			"error", err,
		)
		if i < attempts {
			time.Sleep(2 * time.Second)
		}
	}

	return nil, fmt.Errorf("failed to connect database after %d attempts: %w", attempts, lastErr)
}
