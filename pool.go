package main

import (
	"context"
	"sync"
	"time"

	"github.com/Lz-Gustavo/wormhole/db"
	"github.com/Lz-Gustavo/wormhole/flags"
)

type Pool struct {
	size    int
	prop    flags.Flags
	workers []*Worker

	wg     *sync.WaitGroup
	cancel context.CancelFunc
}

func NewPool(prop flags.Flags) *Pool {
	return &Pool{
		size:    prop.NumClients,
		prop:    prop,
		workers: make([]*Worker, prop.NumClients),
		wg:      &sync.WaitGroup{},
	}
}

// Run ...
func (p *Pool) Run(ctx context.Context, newClient func() db.DatabaseClient) {
	ctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel
	go p.shutdownAfterDur(p.prop.ExecTime)

	for i := range p.size {
		cl := newClient()
		w := NewWorker(cl, p.prop)
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
