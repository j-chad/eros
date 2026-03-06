package storage

import (
	"backend/internal/crypto"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type LocalFileStore struct {
	root string
}

var _ FileStore = &LocalFileStore{}

func NewLocalFileStore(root string) (*LocalFileStore, error) {
	if root == "" {
		return nil, fmt.Errorf("missing required root path")
	}

	if err := os.MkdirAll(root, 0700); err != nil {
		return nil, fmt.Errorf("failed to create root directory: %w", err)
	}

	return &LocalFileStore{root: root}, nil
}

func (s *LocalFileStore) Put(_ context.Context, r io.Reader) (string, error) {
	key := crypto.UUIDV4()
	dst := filepath.Join(s.root, key)

	f, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}
	return key, nil
}

func (s *LocalFileStore) Get(_ context.Context, key string) (io.ReadCloser, error) {
	p, err := s.safePath(key)
	if err != nil {
		return nil, err
	}
	return os.Open(p)
}

func (s *LocalFileStore) Delete(_ context.Context, key string) error {
	p, err := s.safePath(key)
	if err != nil {
		return err
	}
	if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete file: %w", err)
	}
	return nil
}

// safePath resolves the key to an absolute path and ensures it stays within the root.
func (s *LocalFileStore) safePath(key string) (string, error) {
	full := filepath.Join(s.root, filepath.Clean("/"+key))

	// filepath.Clean + Join can still produce something outside root
	// if the resolved path doesn't start with root, reject it
	if !strings.HasPrefix(full, s.root+string(filepath.Separator)) && full != s.root {
		return "", fmt.Errorf("path traversal detected: %s", key)
	}

	return full, nil
}
