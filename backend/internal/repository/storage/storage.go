package storage

import (
	"backend/internal/config"
	"context"
	"fmt"
	"io"
)

type FileStore interface {
	Put(ctx context.Context, r io.Reader) (string, error)
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
}

func NewFileStore(conf config.FileStorageConfig) (FileStore, error) {
	switch {
	case conf.Type == config.FileStorageLocal:
		return NewLocalFileStore(conf.Local.BasePath)
	default:
		return nil, fmt.Errorf("file storage type %s not supported", conf.Type)
	}
}
