package models

import (
	"log/slog"
	"os"
	"strconv"
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
// Lifecycle events (connection established, retry attempts) are logged via
// log/slog as structured records (JSON on most systems), so the app's
// startup diagnostics stay machine-readable.
//
// It retries for a few seconds on startup: when crudin is started via
// docker-compose the Go process and Postgres boot in parallel, so a
// brief "connection refused" is expected. Only if the database stays
// unreachable does the app panic, since without a database the API
// cannot serve its core endpoints.
func ConnectDatabase() {
	dsn := dbDSN()

	const attempts = 5
	var lastErr error

	for i := 1; i <= attempts; i++ {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			// Connection succeeded: migrate the schema, then publish the handle.
			if err := db.AutoMigrate(&Post{}); err != nil {
				panic("failed to migrate database: " + err.Error())
			}
			DB = db
			slog.Info("database connection established")
			return
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

	panic("failed to connect database after " + strconv.Itoa(attempts) + " attempts: " + lastErr.Error())
}
