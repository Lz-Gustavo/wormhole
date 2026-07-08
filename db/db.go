package db

import (
	"context"

	"github.com/Lz-Gustavo/wormhole/flags"
)

// DatabaseFn ...
type DatabaseFn func(flags.Flags) (DatabaseClient, error)

// DatabaseClient ...
type DatabaseClient interface {
	Write(ctx context.Context, key, value string) error
	Close() error
}
