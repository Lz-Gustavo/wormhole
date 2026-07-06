package etcd

import (
	"context"

	"github.com/Lz-Gustavo/wormhole/db"
)

var _ db.DatabaseClient = &EtcdClient{}

type EtcdClient struct{}

func (ec *EtcdClient) Init() error {
	// TODO
	return nil
}

func (ec *EtcdClient) Write(ctx context.Context, key, value []byte) error {
	// TODO
	return nil
}

func (ec *EtcdClient) Close() error {
	// TODO
	return nil
}
