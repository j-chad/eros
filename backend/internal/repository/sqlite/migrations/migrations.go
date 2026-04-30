package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

//go:embed *.sql
var migrationFiles embed.FS

const migrationTable = "migrations"
const legacyMigrationsTable = "node" // if this table exists we assume the 1st migration has been applied

type Migration struct {
	filename string
	Version  int
	Name     string
}

// MigrationStatus represents the status of a single migration: applied or pending.
type MigrationStatus struct {
	Version   int
	Name      string
	Applied   bool
	AppliedAt *time.Time // nil if pending
}

func migrationFromFilename(filename string) (Migration, error) {
	// Expected format: "001_initial_schema.sql"
	withoutExt := strings.TrimSuffix(filename, ".sql")
	if withoutExt == filename {
		return Migration{}, fmt.Errorf("invalid filename format: %s (must end in .sql)", filename)
	}

	parts := strings.SplitN(withoutExt, "_", 2)
	if len(parts) != 2 {
		return Migration{}, fmt.Errorf("invalid filename format: %s", filename)
	}

	version, err := strconv.Atoi(parts[0])
	if err != nil {
		return Migration{}, fmt.Errorf("invalid version in filename %s: %w", filename, err)
	}

	return Migration{
		filename: filename,
		Version:  version,
		Name:     parts[1],
	}, nil
}

func contents(m Migration) (string, error) {
	data, err := migrationFiles.ReadFile(m.filename)
	if err != nil {
		return "", fmt.Errorf("error reading migration file: %w", err)
	}
	return string(data), nil
}

// getCurrentVersion returns the version that the sqlite db has been migrated to.
//
// A version of 0 indicates no migrations have been applied.
//
// If there are gaps in the versions, an error is returned.
func getCurrentVersion(ctx context.Context, db *sql.DB) (int, error) {
	// language=sqlite
	rows, err := db.QueryContext(ctx, "SELECT version FROM migrations ORDER BY version;")
	if err != nil {
		return 0, fmt.Errorf("error querying migrations: %w", err)
	}
	defer rows.Close()

	// ensure versions are sequential
	var version int
	for rows.Next() {
		var nextVersion int
		err = rows.Scan(&nextVersion)
		if err != nil {
			return 0, fmt.Errorf("error scanning row: %w", err)
		}

		if nextVersion != version+1 {
			return 0, fmt.Errorf("unexpected version. expected %d, got %d", nextVersion, version+1)
		}

		version = nextVersion
	}

	err = rows.Err()
	if err != nil {
		return 0, fmt.Errorf("error iterating rows: %w", err)
	}

	return version, nil
}

// bootstrapMigrations creates the migration table if it doesn't already exist.
func bootstrapMigrations(ctx context.Context, db *sql.DB) error {
	exists, err := tableExists(ctx, db, migrationTable)
	if err != nil {
		return fmt.Errorf("error checking for migrations table: %w", err)
	}

	if exists {
		return nil
	}

	// language=sqlite
	_, err = db.ExecContext(ctx, `
		CREATE TABLE migrations (
			version INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("error creating migrations table: %w", err)
	}

	legacyMigrationsExist, err := tableExists(ctx, db, legacyMigrationsTable)
	if err != nil {
		return fmt.Errorf("error checking for legacy migrations: %w", err)
	}

	// legacy migrations cover the first migration
	if legacyMigrationsExist {
		_, err := db.ExecContext(ctx, `
			INSERT INTO migrations (version, name) VALUES (1, 'legacy_migrations');
		`)
		if err != nil {
			return fmt.Errorf("error inserting legacy migrations: %w", err)
		}
	}

	return nil
}

func tableExists(ctx context.Context, db *sql.DB, tableName string) (bool, error) {
	var count int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?;", tableName).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking for migrations table: %w", err)
	}

	return count > 0, nil
}

// getMigrations reads all migration files from the embedded filesystem, parses their version and name, and returns them as a slice of Migrations.
//
// If there are gaps in the versions, an error is returned.
func getMigrations() []Migration {
	entries, err := migrationFiles.ReadDir(".")
	if err != nil {
		panic(fmt.Sprintf("error reading migration files: %v", err))
	}

	migrations := make(map[int]Migration)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		migration, err := migrationFromFilename(entry.Name())
		if err != nil {
			panic(fmt.Sprintf("error parsing migration filename %s: %v", entry.Name(), err))
		}

		if _, exists := migrations[migration.Version]; exists {
			panic(fmt.Sprintf("duplicate migration version %d in file %s and %s", migration.Version, entry.Name(), migrations[migration.Version].filename))
		}

		migrations[migration.Version] = migration
	}

	migrationList := make([]Migration, 0, len(migrations))
	for i := 1; i <= len(migrations); i++ {
		migration, exists := migrations[i]
		if !exists {
			panic(fmt.Sprintf("missing migration version %d", i))
		}
		migrationList = append(migrationList, migration)
	}

	return migrationList
}

// Remaining returns all migrations that have not yet been applied, in order of their version.
func Remaining(ctx context.Context, db *sql.DB) ([]Migration, error) {
	err := bootstrapMigrations(ctx, db)
	if err != nil {
		return nil, err
	}

	migrations := getMigrations()
	remaining := make([]Migration, 0, len(migrations))

	version, err := getCurrentVersion(ctx, db)
	if err != nil {
		return nil, err
	}

	for _, migration := range migrations {
		if migration.Version > version {
			remaining = append(remaining, migration)
		}
	}

	return remaining, nil
}

// Apply runs all pending migrations in order. Each migration executes inside its
// own transaction. This works because go-sqlite3 uses sqlite3_exec() under the
// hood, which supports multiple statements in a single call.
func Apply(ctx context.Context, db *sql.DB) error {
	remaining, err := Remaining(ctx, db)
	if err != nil {
		return fmt.Errorf("error determining pending migrations: %w", err)
	}

	if len(remaining) == 0 {
		slog.Info("no pending migrations")
		return nil
	}

	for _, m := range remaining {
		sqlContent, err := contents(m)
		if err != nil {
			return fmt.Errorf("error reading migration %03d_%s: %w", m.Version, m.Name, err)
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("error beginning transaction for migration %03d_%s: %w", m.Version, m.Name, err)
		}

		if _, err := tx.ExecContext(ctx, sqlContent); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("error executing migration %03d_%s: %w", m.Version, m.Name, err)
		}

		// language=sqlite
		if _, err := tx.ExecContext(ctx, "INSERT INTO migrations (version, name) VALUES (?, ?);", m.Version, m.Name); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("error recording migration %03d_%s: %w", m.Version, m.Name, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("error committing migration %03d_%s: %w", m.Version, m.Name, err)
		}

		slog.Info("applied migration", "migration", fmt.Sprintf("%03d_%s", m.Version, m.Name))
	}

	return nil
}

// Status returns the status of all known migrations: applied (with timestamp) or pending.
func Status(ctx context.Context, db *sql.DB) ([]MigrationStatus, error) {
	err := bootstrapMigrations(ctx, db)
	if err != nil {
		return nil, err
	}

	allMigrations := getMigrations()

	// Read applied migrations with their timestamps.
	applied := make(map[int]time.Time)
	// language=sqlite
	rows, err := db.QueryContext(ctx, "SELECT version, applied_at FROM migrations ORDER BY version;")
	if err != nil {
		return nil, fmt.Errorf("error querying migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var version int
		var appliedAt time.Time
		if err := rows.Scan(&version, &appliedAt); err != nil {
			return nil, fmt.Errorf("error scanning migration row: %w", err)
		}
		applied[version] = appliedAt
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating migration rows: %w", err)
	}

	statuses := make([]MigrationStatus, 0, len(allMigrations))
	for _, m := range allMigrations {
		status := MigrationStatus{
			Version: m.Version,
			Name:    m.Name,
		}
		if at, ok := applied[m.Version]; ok {
			status.Applied = true
			status.AppliedAt = &at
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}
