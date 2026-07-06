package db

import "context"

// DatabaseClient ...
type DatabaseClient interface {
	Init() error
	Write(ctx context.Context, key, value []byte) error
	Close() error
}
