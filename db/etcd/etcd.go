package etcd

import (
	"context"
	"math/rand"
	"time"

	"github.com/Lz-Gustavo/wormhole/db"
	"github.com/Lz-Gustavo/wormhole/flags"
	"github.com/Lz-Gustavo/wormhole/measure"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

const (
	etcdDialTimeout = 5 * time.Second

	// each request has a '1/measureChance' chance to record latency and/or
	// the response status code
	measureChance = 10
)

var (
	_ db.DatabaseClient = &EtcdClient{}
	_ db.NewDatabaseFn  = NewEtcdClient
)

type EtcdClient struct {
	meter *measure.Meter
	prop  flags.Flags

	client *clientv3.Client
}

func NewEtcdClient(prop flags.Flags, meter *measure.Meter) (db.DatabaseClient, error) {
	cfg := clientv3.Config{
		Endpoints:   prop.EtcdHosts,
		DialTimeout: etcdDialTimeout,
	}

	if !prop.Verbose {
		cfg.Logger = zap.NewNop()
	}

	cl, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}

	return &EtcdClient{
		client: cl,
		prop:   prop,
		meter:  meter,
	}, nil
}

func (ec *EtcdClient) Write(ctx context.Context, key, value string) error {
	var err error

	if ec.isLatencyMsrEnabled() && mustMeasureLat() {
		start := time.Now()
		_, err = ec.client.Put(ctx, key, value)
		if errL := ec.meter.LatMsr.Record(time.Since(start)); errL != nil {
			return errL
		}

	} else {
		_, err = ec.client.Put(ctx, key, value)
	}

	if ec.isStatusMsrEnabled() {
		ec.meter.StatusMsr.CountStatusFromErr(err)
	}
	return err
}

func (ec *EtcdClient) Close() error {
	if ec.isLatencyMsrEnabled() {
		if err := ec.meter.LatMsr.Flush(); err != nil {
			return err
		}
		if err := ec.meter.LatMsr.Close(); err != nil {
			return err
		}
	}

	if ec.isStatusMsrEnabled() {
		if err := ec.meter.StatusMsr.Flush(); err != nil {
			return err
		}
		if err := ec.meter.StatusMsr.Close(); err != nil {
			return err
		}
	}
	return ec.client.Close()
}

func (ec *EtcdClient) isLatencyMsrEnabled() bool {
	return ec.meter != nil && ec.meter.LatMsr != nil
}

func (ec *EtcdClient) isStatusMsrEnabled() bool {
	return ec.meter != nil && ec.meter.StatusMsr != nil
}

func mustMeasureLat() bool {
	return rand.Intn(measureChance) == 0
}
