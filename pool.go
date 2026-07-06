package main

import (
	"context"
	"sync"
	"time"

	"github.com/Lz-Gustavo/wormhole/db"
)

type Pool struct {
	size    int
	workers []*Worker

	wg     *sync.WaitGroup
	cancel context.CancelFunc
}

func NewPool(size int) *Pool {
	return &Pool{
		size:    size,
		workers: make([]*Worker, size),
		wg:      &sync.WaitGroup{},
	}
}

// Run ...
func (p *Pool) Run(ctx context.Context, newClient func() db.DatabaseClient, dur time.Duration) {
	ctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel
	go p.shutdownAfterDur(dur)

	for i := range p.size {
		cl := newClient()
		w := NewWorker(cl)
		p.workers[i] = w

		p.wg.Go(func() {
			w.Run(ctx)
		})
	}
	p.wg.Wait()
}

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

	p.cancel()
}
