package testdb

import (
	"backend/internal/config"
	"backend/internal/repository"
	"backend/internal/repository/sqlite"
	"backend/internal/repository/storage"
	"testing"
)

// New creates an in-memory SQLite repository for testing.
// Uses cache=shared so all connections in the pool share the same database.
// Automatically closed when the test finishes.
func New(t *testing.T) repository.Repository {
	t.Helper()
	repo, err := sqlite.NewSQLiteDB(config.DatabaseConfig{
		Path: "file::memory:?cache=shared",
	})
	if err != nil {
		t.Fatalf("failed to create test repo: %v", err)
	}
	t.Cleanup(func() { repo.Close() })
	return repo
}

// NewFileStore creates a local file store backed by a temp directory.
// Automatically cleaned up when the test finishes.
func NewFileStore(t *testing.T) *storage.LocalFileStore {
	t.Helper()
	store, err := storage.NewLocalFileStore(t.TempDir())
	if err != nil {
		t.Fatalf("failed to create test file store: %v", err)
	}
	return store
}
