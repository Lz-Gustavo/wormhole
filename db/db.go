package db

import (
	"context"

	"github.com/Lz-Gustavo/wormhole/flags"
	"github.com/Lz-Gustavo/wormhole/measure"
)

// NewDatabaseFn is a common signature for database constructors to be later switched at runtime.
type NewDatabaseFn func(flags.Flags, *measure.Meter) (DatabaseClient, error)

// DatabaseClient defines the database interface utilized by workers.
type DatabaseClient interface {
	Write(ctx context.Context, key, value string) error
	Close() error
}
