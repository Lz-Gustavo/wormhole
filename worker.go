package main

import (
	"context"
	"log"
	"time"

	"github.com/Lz-Gustavo/wormhole/db"
)

const (
	defaultCommandTimeout = 5 * time.Second

	defaultKeySpaceSize  = 10000
	defaultPayloadSizeKb = 4
)

type Worker struct {
	client db.DatabaseClient
	logger *log.Logger

	thinkingTime time.Duration
}

func NewWorker(cl db.DatabaseClient, tt time.Duration) *Worker {
	return &Worker{
		logger:       log.Default(),
		client:       cl,
		thinkingTime: tt,
	}
}

func (w *Worker) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		default:
			w.work(ctx)
		}
	}
}

func (w *Worker) work(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, defaultCommandTimeout)
	defer cancel()

	key := db.GetRandKeyUpTo(defaultKeySpaceSize)
	value := db.GetPayloadBySizeKb(defaultPayloadSizeKb)

	if err := w.client.Write(ctx, key, value); err != nil {
		w.logger.Println("failed write request: %w", err)
	}
	time.Sleep(w.getThinkingTime())
}

func (w *Worker) getThinkingTime() time.Duration {
	// TODO: generate rand from 0 to w.thinkingTime
	return time.Second
}
