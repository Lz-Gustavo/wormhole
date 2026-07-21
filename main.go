package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/Lz-Gustavo/wormhole/db/etcd"
	"github.com/Lz-Gustavo/wormhole/flags"
)

func getLogHandler(f flags.Flags) slog.Handler {
	level := slog.LevelInfo
	if f.Verbose {
		level = slog.LevelDebug
	}

	return slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format(time.DateTime))
				}
			}
			return a
		},
	})
}

func main() {
	f := flags.ParseFlagsFromArgs()
	slog.SetDefault(slog.New(getLogHandler(f)).With("src", "wormhole"))

	ctx := context.Background()

	p, err := NewPool(f)
	if err != nil {
		slog.Error("failed initializing worker pool", "err", err)
		os.Exit(1)
	}

	p.Run(ctx, etcd.NewEtcdClient)
	slog.Info("finished", "success-req-count", p.Count())
}
