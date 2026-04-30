package sqlite

import (
	"backend/internal/config"
	"backend/internal/repository"
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

//goland:noinspection GoResourceLeak
func OpenDB(conf config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", conf.Path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if conf.WAL {
		if _, err := db.Exec("PRAGMA journal_mode = WAL;"); err != nil {
			return nil, fmt.Errorf("error setting WAL mode: %w", err)
		}
	}

	if conf.BusyTimeout > 0 {
		busyTimeoutMS := int(conf.BusyTimeout.Milliseconds())
		query := fmt.Sprintf("PRAGMA busy_timeout = %d;", busyTimeoutMS)
		if _, err := db.Exec(query); err != nil {
			return nil, fmt.Errorf("error setting busy timeout: %w", err)
		}
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return nil, fmt.Errorf("error enabling foreign keys: %w", err)
	}

	return db, nil
}

type Queryable interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// NewSQLiteDB opens a SQLite connection and wraps it as a repository.
// It does NOT run migrations — the caller must run migrations separately
// via the migrate command. If the database hasn't been migrated, the server
// will fail on first query, making it obvious that migrations need to be run.
func NewSQLiteDB(conf config.DatabaseConfig) (repository.Repository, error) {
	db, err := OpenDB(conf)
	if err != nil {
		return nil, err
	}

	return NewFromDB(db), nil
}

// NewFromDB wraps a pre-opened *sql.DB as a repository.Repository.
// Does not configure PRAGMAs or run schema init — the caller is
// responsible for that (via OpenDB and/or the migration runner).
func NewFromDB(db *sql.DB) repository.Repository {
	return &sqliteDB{db: db}
}

type sqliteDB struct {
	db *sql.DB // DEPRECATED: use executor() instead to support transactions
	tx *sql.Tx
}

func (s *sqliteDB) Close() {
	if err := s.db.Close(); err != nil {
		fmt.Printf("error closing database: %v\n", err)
	}
}

func (s *sqliteDB) WithTx(ctx context.Context, opts *sql.TxOptions, fn func(repository.Repository) error) error {
	return s.withTx(ctx, opts, func(db *sqliteDB) error {
		return fn(db)
	})
}

func (s *sqliteDB) withTx(ctx context.Context, opts *sql.TxOptions, fn func(db *sqliteDB) error) error {
	if s.tx != nil {
		// already in a transaction
		return fn(s)
	}

	tx, err := s.db.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	txRepo := &sqliteDB{
		db: s.db,
		tx: tx,
	}

	if err := fn(txRepo); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("error during transaction rollback: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	if opts != nil && opts.ReadOnly {
		// read-only transaction, no need to commit
		return nil
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func (s *sqliteDB) executor() Queryable {
	if s.tx != nil {
		return s.tx
	}

	return s.db
}
