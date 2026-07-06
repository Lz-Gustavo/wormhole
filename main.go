package main

import (
	"context"
	"time"

	"github.com/Lz-Gustavo/wormhole/db/etcd"
)

const (
	defaultSize     = 10
	defaultDuration = 5 * time.Second
)

func main() {
	ctx := context.Background()
	p := NewPool(defaultSize)

	p.Run(ctx, etcd.NewEtcdClient, defaultDuration)
}
