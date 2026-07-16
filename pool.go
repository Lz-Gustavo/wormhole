package main

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/Lz-Gustavo/wormhole/db"
	"github.com/Lz-Gustavo/wormhole/flags"
)

// Pool coordinates workers lifecycle, starts one database client per worker, and gracefully
// stops them after the configured execution time.
type Pool struct {
	size    int
	prop    flags.Flags
	workers []*Worker
	logger  *slog.Logger

	wg     *sync.WaitGroup
	cancel context.CancelFunc
}

// NewPool builds a Pool with the requested number of clients and a default logger.
func NewPool(prop flags.Flags) *Pool {
	return &Pool{
		size:    prop.NumClients,
		prop:    prop,
		workers: make([]*Worker, prop.NumClients),
		logger:  slog.Default(),
		wg:      &sync.WaitGroup{},
	}
}

// Run creates the database clients, starts each worker, and cancels them after ExecTime.
func (p *Pool) Run(ctx context.Context, newClient db.NewDatabaseFn) {
	ctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel
	go p.shutdownAfterDur(p.prop.ExecTime)

	for i := range p.size {
		cl, err := newClient(p.prop)
		if err != nil {
			p.logger.Error("failed initializing database client", "err", err)
			return
		}

		w := NewWorker(cl, p.prop)
		p.workers[i] = w

		p.wg.Go(func() {
			w.Run(ctx)
		})
	}
	p.wg.Wait()
}

// Count returns the total number of successful requests reported by all workers.
func (p *Pool) Count() int64 {
	var n int64
	for _, w := range p.workers {
		n += w.Count()
	}
	return n
}

func (p *Pool) shutdownAfterDur(dur time.Duration) {
	t := time.NewTimer(dur)
	<-t.C

	p.logger.Debug("shutting down workers...")
	p.cancel()
}
