package storage

import (
	"backend/internal/config"
	"context"
	"fmt"
	"io"
	"iter"
	"time"
)

// PresignCapable Optional interface that storage backends can implement
type PresignCapable interface {
	PresignGet(ctx context.Context, key string, ttl time.Duration) (string, error)
}

type FileStore interface {
	Put(ctx context.Context, filename string, r io.ReadSeeker) (string, error)
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	List(ctx context.Context) iter.Seq2[string, error]
	DeleteMany(ctx context.Context, keys []string) error
}

func NewFileStore(conf config.FileStorageConfig) (FileStore, error) {
	switch {
	case conf.Type == config.FileStorageLocal:
		return NewLocalFileStore(conf.Local.BasePath)
	case conf.Type == config.FileStorageS3:
		return NewS3FileStore(conf.S3), nil
	default:
		return nil, fmt.Errorf("file storage type %s not supported", conf.Type)
	}
}
