package logger

import (
	"sync/atomic"
	"time"
)

type ticker struct {
	ticker <-chan time.Time
	// collection of &Logger.count
	counters []*int64
}

func (t *ticker) register(counter *int64) {
	t.counters = append(t.counters, counter)
}

func newTicker() *ticker {
	t := &ticker{
		ticker: time.Tick(time.Second),
	}
	go func() {
		for {
			<-t.ticker
			for _, v := range t.counters {
				atomic.StoreInt64(v, 0)
			}
		}
	}()
	return t
}

var globalTicker = newTicker()
