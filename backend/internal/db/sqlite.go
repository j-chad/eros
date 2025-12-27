package db

import (
	"database/sql"
	_ "embed"
	"fmt"
)

//go:embed sqlite_init.sql
var sqliteInitSQL string

func newSQLiteDB() (Database, error) {
	//goland:noinspection GoResourceLeak connection is designed to be long-lived
	db, err := sql.Open("sqlite", "./db/sqlite.db")
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
}

func (s sqliteDB) Close() {
	if err := s.db.Close(); err != nil {
		fmt.Printf("error closing database: %v\n", err)
	}
}
