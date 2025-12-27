package sqlite

import (
	"backend/internal/repository"
	"context"
	"database/sql"
	_ "embed"
	"fmt"
)

//go:embed sqlite_init.sql
var sqliteInitSQL string

type executor interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

func NewSQLiteDB() (repository.Repository, error) {
	//goland:noinspection GoResourceLeak connection is designed to be long-lived
	db, err := sql.Open("sqlite", "./repository.sqlite")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// initialize the database schema
	if _, err := db.Exec(sqliteInitSQL); err != nil {
		return nil, fmt.Errorf("error initializing database: %w", err)
	}

	return &sqliteDB{db: db}, nil
}

type sqliteDB struct {
	db *sql.DB
	tx *sql.Tx
}

func (s sqliteDB) Close() {
	if err := s.db.Close(); err != nil {
		fmt.Printf("error closing database: %v\n", err)
	}
}

func (s sqliteDB) WithTx(ctx context.Context, fn func(repository.Repository) error) error {
	if s.tx != nil {
		return fmt.Errorf("nested transactions are not supported")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

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

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func (s sqliteDB) executor() executor {
	if s.tx != nil {
		return s.tx
	}

	return s.db
}
