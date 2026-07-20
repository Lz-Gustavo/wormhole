package main

import (
	"context"
	"errors"
	"log/slog"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Lz-Gustavo/wormhole/db"
	"github.com/Lz-Gustavo/wormhole/flags"
)

// Worker owns a database client and tracks its request count.
type Worker struct {
	client db.DatabaseClient
	prop   flags.Flags
	logger *slog.Logger

	count atomic.Int64
}

// NewWorker instantiates a Worker with the provided client and configuration.
func NewWorker(cl db.DatabaseClient, prop flags.Flags) *Worker {
	return &Worker{
		prop:   prop,
		logger: slog.Default(),
		client: cl,
	}
}

// Run implements an open-loop iteration, where every request is assigned to background goroutine
// and a new one is issued every random thinking time.
func (w *Worker) Run(ctx context.Context) {
	w.logger.Debug("worker started...")
	var wg sync.WaitGroup

	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			if err := w.client.Close(); err != nil {
				w.logger.Error("failed closing database client", "err", err)
			}
			w.logger.Debug("worker finished")
			return

		default:
			wg.Go(func() {
				w.work(ctx)
			})
			time.Sleep(w.getRandThinkingTime())
		}
	}
}

// Count returns the number of completed requests by the worker.
func (w *Worker) Count() int64 {
	return w.count.Load()
}

func (w *Worker) work(ctx context.Context) {
	key := db.GetRandKeyUpTo(w.prop.KeySpaceSize)
	value := db.GetPayloadBySizeKb(db.PayloadSize(w.prop.PayloadSize))

	ctx, cancel := context.WithTimeout(ctx, w.prop.CmdTimeout)
	defer cancel()

	if err := w.client.Write(ctx, key, value); err != nil {
		if w.shouldLog(err) {
			w.logger.Error("failed on write request", "err", err)
		}
		return
	}
	w.count.Add(1)
}

func (w *Worker) getRandThinkingTime() time.Duration {
	ms := rand.IntN(w.prop.MaxThinkingTimeMs)
	return time.Duration(ms) * time.Millisecond
}

func (w *Worker) shouldLog(err error) bool {
	return w.prop.Verbose ||
		!errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded)
}
