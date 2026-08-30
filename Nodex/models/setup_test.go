package models

import "testing"

// TestDBDSNEnvOverride confirms DATABASE_URL takes precedence over the local
// default — so the app is configurable entirely via the environment.
func TestDBDSNEnvOverride(t *testing.T) {
	const custom = "host=10.0.0.9 user=foo dbname=bar sslmode=disable"
	t.Setenv("DATABASE_URL", custom)

	if got := dbDSN(); got != custom {
		t.Fatalf("dbDSN() = %q, want %q (DATABASE_URL override ignored)", got, custom)
	}
}

// TestDBDSNDefaults confirms the local default DSN is used when DATABASE_URL is
// unset, so the app runs with zero config on the practice setup.
func TestDBDSNDefaults(t *testing.T) {
	t.Setenv("DATABASE_URL", "")

	if got := dbDSN(); got != defaultDSN {
		t.Fatalf("dbDSN() = %q, want default %q", got, defaultDSN)
	}
}
