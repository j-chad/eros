package storage

import (
	"backend/internal/crypto"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

var ALLOWED_KEY_REGEX = regexp.MustCompile(fmt.Sprintf("^%s(\\.[a-zA-Z]{1,3})?$", crypto.UUIDV4Regex.String()))

type LocalFileStore struct {
	root string
}

func NewLocalFileStore(root string) (*LocalFileStore, error) {
	if root == "" {
		return nil, fmt.Errorf("missing required root path")
	}

	if err := os.MkdirAll(root, 0700); err != nil {
		return nil, fmt.Errorf("failed to create root directory: %w", err)
	}

	return &LocalFileStore{root: root}, nil
}

func (s *LocalFileStore) Put(_ context.Context, filename string, r io.Reader) (string, error) {
	key := s.getKey(filename)
	dst, err := s.safePath(key)
	if err != nil {
		return "", err
	}

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
	if !ALLOWED_KEY_REGEX.MatchString(key) {
		return "", fmt.Errorf("invalid key format")
	}

	return filepath.Join(s.root, key), nil
}

func (s *LocalFileStore) getKey(filename string) string {
	ext := filepath.Ext(filename)
	uuid := crypto.UUIDV4()
	return uuid + ext
}
