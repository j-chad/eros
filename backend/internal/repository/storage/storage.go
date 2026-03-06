package storage

import (
	"context"
	"io"
)

type FileStore interface {
	Put(ctx context.Context, r io.Reader) (string, error)
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
}
