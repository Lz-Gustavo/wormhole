package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Lz-Gustavo/wormhole/db/etcd"
	"github.com/Lz-Gustavo/wormhole/flags"
)

func main() {
	f := flags.ParseFlagsFromArgs()
	level := slog.LevelError
	if f.Verbose {
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: level})))

	ctx := context.Background()

	p := NewPool(f)
	p.Run(ctx, etcd.NewEtcdClient)
	slog.Info("finished", "count", p.Count())
}
