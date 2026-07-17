package main

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Lz-Gustavo/wormhole/db"
	"github.com/Lz-Gustavo/wormhole/flags"
	"github.com/stretchr/testify/assert"
)

type workerTestClient struct {
	writes atomic.Int64
}

func (c *workerTestClient) Write(ctx context.Context, key, value string) error {
	c.writes.Add(1)
	return nil
}

func (c *workerTestClient) Close() error {
	return nil
}

func Test_WorkerRun(t *testing.T) {
	testExecTime := 100 * time.Millisecond

	cl := &workerTestClient{}
	w := NewWorker(cl, flags.Flags{
		CmdTimeout:        time.Millisecond,
		MaxThinkingTimeMs: 1,
		KeySpaceSize:      10,
		PayloadSize:       int(db.Small),
	})

	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan struct{})
	go func() {
		w.Run(ctx)
		close(done)
	}()

	time.Sleep(testExecTime)
	cancel()

	select {
	case <-done:
	case <-time.After(testExecTime):
		t.Fatal("worker did not stop after cancellation")
	}

	assert.Greater(t, w.Count(), int64(0))
	assert.Equal(t, cl.writes.Load(), w.Count())
}
