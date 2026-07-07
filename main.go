package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Lz-Gustavo/wormhole/db/etcd"
	"github.com/Lz-Gustavo/wormhole/flags"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	f := flags.ParseFlagsFromArgs()
	ctx := context.Background()

	p := NewPool(f)
	p.Run(ctx, etcd.NewEtcdClient)
	fmt.Println("executed:", p.Count())
}
