package logger

import (
	"sync/atomic"
	"time"
)

type ticker struct {
	ticker  <-chan time.Time
	loggers []*Logger
}

func (t *ticker) register(logger *Logger) {
	t.loggers = append(t.loggers, logger)
}

func newTicker() *ticker {
	t := &ticker{
		ticker: time.Tick(time.Second),
	}
	go func() {
		for {
			<-t.ticker
			for _, v := range t.loggers {
				atomic.StoreInt64(v.count, 0)
			}
		}
	}()
	return t
}

var globalTicker = newTicker()
