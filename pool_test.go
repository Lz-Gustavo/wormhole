package main

import (
	"context"
	"testing"
	"time"

	"github.com/Lz-Gustavo/wormhole/db"
	"github.com/Lz-Gustavo/wormhole/flags"
	"github.com/stretchr/testify/assert"
)

func Test_PoolRun(t *testing.T) {
	testExecTime := 200 * time.Millisecond
	waitFor := 5 * testExecTime

	prop := flags.Flags{
		NumClients:        4,
		ExecTime:          testExecTime,
		CmdTimeout:        time.Millisecond,
		MaxThinkingTimeMs: 10,
		KeySpaceSize:      10,
		PayloadSize:       int(db.Small),
	}

	p := NewPool(prop)
	newClient := func(flags.Flags) (db.DatabaseClient, error) {
		return &workerTestClient{}, nil
	}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	done := make(chan struct{})
	go func() {
		p.Run(ctx, newClient)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(waitFor):
		t.Fatal("pool did not stop")
	}

	assert.GreaterOrEqual(t, p.Count(), int64(prop.NumClients))
}
