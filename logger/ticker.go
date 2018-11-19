package logger

import (
	"sync"
	"sync/atomic"
	"time"
)

type ticker struct {
	mut    sync.Mutex
	ticker <-chan time.Time
	// collection of &Logger.count

	counters []*int64
}

func (t *ticker) register(counter *int64) {
	t.mut.Lock()
	t.counters = append(t.counters, counter)
	t.mut.Unlock()
}

func newTicker() *ticker {
	t := &ticker{
		ticker: time.Tick(time.Second),
	}
	go func() {
		for {
			<-t.ticker
			t.mut.Lock()
			for _, v := range t.counters {
				atomic.StoreInt64(v, 0)
			}
			t.mut.Unlock()
		}
	}()
	return t
}

var globalTicker = newTicker()
