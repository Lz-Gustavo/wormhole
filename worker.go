package main

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"sync/atomic"
	"time"

	"github.com/Lz-Gustavo/wormhole/db"
	"github.com/Lz-Gustavo/wormhole/flags"
)

type Worker struct {
	client db.DatabaseClient
	prop   flags.Flags
	logger *slog.Logger

	count atomic.Int64
}

func NewWorker(cl db.DatabaseClient, prop flags.Flags) *Worker {
	return &Worker{
		prop:   prop,
		logger: slog.Default(),
		client: cl,
	}
}

func (w *Worker) Run(ctx context.Context) {
	w.logger.Debug("worker started...")

	for {
		select {
		case <-ctx.Done():
			if err := w.client.Close(); err != nil {
				w.logger.Error("failed closing database client", "err", err)
			}
			w.logger.Debug("worker finished")
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
	key := db.GetRandKeyUpTo(w.prop.KeySpaceSize)
	value := db.GetPayloadBySizeKb(db.PayloadSize(w.prop.PayloadSize))

	ctx, cancel := context.WithTimeout(ctx, w.prop.CmdTimeout)
	defer cancel()

	if err := w.client.Write(ctx, key, value); err != nil {
		w.logger.Error("failed on write request", "err", err)
	}
	w.count.Add(1)
}

func (w *Worker) getRandThinkingTime() time.Duration {
	ms := rand.IntN(w.prop.MaxThinkingTimeMs)
	return time.Duration(ms) * time.Millisecond
}
