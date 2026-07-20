package etcd

import (
	"context"
	"math/rand"
	"time"

	"github.com/Lz-Gustavo/wormhole/db"
	"github.com/Lz-Gustavo/wormhole/flags"
	"github.com/Lz-Gustavo/wormhole/measure"
	clientv3 "go.etcd.io/etcd/client/v3"
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
	client *clientv3.Client
	prop   flags.Flags

	isLatencyMsrEnabled bool
	latMsr              *measure.LatencyMsr

	isStatusMsrEnabled bool
	statusMsr          *measure.StatusMsr
	stopMsr            context.CancelFunc
}

func NewEtcdClient(prop flags.Flags) (db.DatabaseClient, error) {
	cl, err := clientv3.New(clientv3.Config{
		Endpoints:   prop.EtcdHosts,
		DialTimeout: etcdDialTimeout,
	})
	if err != nil {
		return nil, err
	}

	ec := &EtcdClient{
		client: cl,
		prop:   prop,
	}

	if err := ec.initializeMsrFromProp(); err != nil {
		return nil, err
	}
	return ec, nil
}

func (ec *EtcdClient) Write(ctx context.Context, key, value string) error {
	var err error

	if ec.isLatencyMsrEnabled && mustMeasureLat() {
		start := time.Now()
		_, err = ec.client.Put(ctx, key, value)
		if errL := ec.latMsr.Record(time.Since(start)); errL != nil {
			return errL
		}

	} else {
		_, err = ec.client.Put(ctx, key, value)
	}

	if ec.isStatusMsrEnabled {
		ec.statusMsr.CountStatusFromErr(err)
	}
	return err
}

func (ec *EtcdClient) Close() error {
	if ec.isLatencyMsrEnabled {
		if err := ec.latMsr.Flush(); err != nil {
			return err
		}
		if err := ec.latMsr.Close(); err != nil {
			return err
		}
	}

	if ec.isStatusMsrEnabled {
		ec.stopMsr()
		if err := ec.statusMsr.Flush(); err != nil {
			return err
		}
		if err := ec.statusMsr.Close(); err != nil {
			return err
		}
	}
	return ec.client.Close()
}

func (ec *EtcdClient) initializeMsrFromProp() error {
	if ec.prop.LatencyMsrFilename != "" {
		lm, err := measure.NewLatencyMsr(ec.prop.LatencyMsrFilename)
		if err != nil {
			return err
		}
		ec.isLatencyMsrEnabled = true
		ec.latMsr = lm
	}

	if ec.prop.StatusMsrFilename != "" {
		sm, err := measure.NewStatusMsr(ec.prop.StatusMsrFilename)
		if err != nil {
			return err
		}
		ec.isStatusMsrEnabled = true
		ec.statusMsr = sm

		ctx, cancel := context.WithCancel(context.Background())
		ec.stopMsr = cancel
		go ec.statusMsr.Run(ctx)
	}
	return nil
}

func mustMeasureLat() bool {
	return rand.Intn(measureChance) == 0
}
