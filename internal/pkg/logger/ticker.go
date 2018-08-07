package logger

import (
	"sync"
	"time"
)

type ticker struct {
	mu      sync.Mutex
	clients []chan struct{}
}

func (g *ticker) register(c chan struct{}) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.clients = append(g.clients, c)
}

func (g *ticker) notify() {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, c := range g.clients {
		select {
		case c <- struct{}{}:
		default:
		}
	}
}

func newTicker() *ticker {
	o := new(ticker)
	t := time.Tick(time.Second)
	go func() {
		for range t {
			o.notify()
		}
	}()
	return o
}

var globalTicker = newTicker()
