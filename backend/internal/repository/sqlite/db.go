package sqlite

import (
	"backend/internal/config"
	"backend/internal/repository"
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed sqlite_init.sql
var sqliteInitSQL string

type executor interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func NewSQLiteDB(conf config.DatabaseConfig) (repository.Repository, error) {
	//goland:noinspection GoResourceLeak connection is designed to be long-lived
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

	// initialize the database schema
	if _, err := db.Exec(sqliteInitSQL); err != nil {
		return nil, fmt.Errorf("error initializing database: %w", err)
	}

	return &sqliteDB{db: db}, nil
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

func (s *sqliteDB) executor() executor {
	if s.tx != nil {
		return s.tx
	}

	return s.db
}
