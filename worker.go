package main

import (
	"context"
	"log"
	"math/rand/v2"
	"sync/atomic"
	"time"

	"github.com/Lz-Gustavo/wormhole/db"
)

const (
	defaultCommandTimeout    = 5 * time.Second
	defaultMaxThinkingTimeMs = 500

	defaultKeySpaceSize  = 10000
	defaultPayloadSizeKb = 4
)

type Worker struct {
	client db.DatabaseClient
	logger *log.Logger

	count atomic.Int64
}

func NewWorker(cl db.DatabaseClient) *Worker {
	return &Worker{
		logger: log.Default(),
		client: cl,
	}
}

func (w *Worker) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		default:
			go w.work(ctx)
			time.Sleep(w.getRandThinkingTime())
		}
	}
}

func (w *Worker) Count() int64 {
	return w.count.Load()
}

func (w *Worker) work(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, defaultCommandTimeout)
	defer cancel()

	key := db.GetRandKeyUpTo(defaultKeySpaceSize)
	value := db.GetPayloadBySizeKb(defaultPayloadSizeKb)

	if err := w.client.Write(ctx, key, value); err != nil {
		w.logger.Println("failed write request: %w", err)
	}
	w.count.Add(1)
}

func (w *Worker) getRandThinkingTime() time.Duration {
	ms := rand.IntN(defaultMaxThinkingTimeMs)
	return time.Duration(ms) * time.Millisecond
}
