package measure

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

const (
	initialStatusBuffCap = 1 << 11

	statusFmt = "200:%d-408:%d-500:%d\n"
)

type StatusMsr struct {
	tick         *time.Ticker
	countSuccess atomic.Uint32
	countTimeout atomic.Uint32
	countFail    atomic.Uint32

	mu      sync.Mutex
	buff    *bytes.Buffer
	file    *os.File
	stopMsr context.CancelFunc
}

func NewStatusMsr(filename string) (*StatusMsr, error) {
	sm := &StatusMsr{
		tick: time.NewTicker(time.Second),
		buff: &bytes.Buffer{},
	}
	sm.buff.Grow(initialStatusBuffCap)

	fd, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}
	sm.file = fd
	return sm, nil
}

func (sm *StatusMsr) CountStatusFromErr(err error) {
	if err == nil {
		sm.countSuccess.Add(1)

	} else if errors.Is(err, context.DeadlineExceeded) {
		sm.countTimeout.Add(1)

	} else {
		sm.countFail.Add(1)
	}
}

func (sm *StatusMsr) Run(ctx context.Context) {
	ctx, sm.stopMsr = context.WithCancel(ctx)

	for {
		select {
		case <-ctx.Done():
			return

		case <-sm.tick.C:
			countS := sm.countSuccess.Swap(0)
			countT := sm.countTimeout.Swap(0)
			countF := sm.countFail.Swap(0)

			sm.mu.Lock()
			if _, err := fmt.Fprintf(sm.buff, statusFmt, countS, countT, countF); err != nil {
				sm.mu.Unlock()
				log.Fatalln("failed writing status measurement, err:", err.Error())
				return
			}
			sm.mu.Unlock()
		}
	}
}

func (sm *StatusMsr) Flush() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, err := sm.buff.WriteTo(sm.file); err != nil {
		return err
	}

	if err := sm.file.Sync(); err != nil {
		return err
	}
	return nil
}

func (sm *StatusMsr) Close() error {
	sm.tick.Stop()
	sm.stopMsr()
	return sm.file.Close()
}
