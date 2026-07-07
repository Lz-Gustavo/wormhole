package main

import (
	"context"
	"fmt"

	"github.com/Lz-Gustavo/wormhole/db/etcd"
	"github.com/Lz-Gustavo/wormhole/flags"
)

func main() {
	f := flags.ParseFlagsFromArgs()
	ctx := context.Background()

	p := NewPool(f.NumClients)
	p.Run(ctx, etcd.NewEtcdClient, f.ExecTime)
	fmt.Println("executed:", p.Count())
}
