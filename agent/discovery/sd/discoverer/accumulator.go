// SPDX-License-Identifier: GPL-3.0-or-later

package discoverer

import (
	"context"
	"sync"
	"time"

	"github.com/netdata/go.d.plugin/agent/discovery/sd/model"
)

func newAccumulator() *accumulator {
	return &accumulator{
		send:      make(chan struct{}, 1),
		sendEvery: time.Second * 3,
		mux:       &sync.Mutex{},
		groups:    make(map[string]model.TargetGroup),
	}
}

type accumulator struct {
	discoverers []Discoverer
	send        chan struct{}
	sendEvery   time.Duration
	mux         *sync.Mutex
	groups      map[string]model.TargetGroup
}

func (a *accumulator) run(ctx context.Context, in chan<- []model.TargetGroup) {
	var wg sync.WaitGroup

	for _, d := range a.discoverers {
		wg.Add(1)
		go func(d Discoverer) { defer wg.Done(); a.runDiscoverer(ctx, d, in) }(d)
	}

	wg.Wait()
	<-ctx.Done()
}

func (a *accumulator) runDiscoverer(ctx context.Context, d Discoverer, in chan<- []model.TargetGroup) {
	updates := make(chan []model.TargetGroup)
	go d.Discover(ctx, updates)

	for {
		select {
		case <-ctx.Done():
			return
		case groups, ok := <-updates:
			if !ok {
				return
			}
			a.mux.Lock()
			a.groupsUpdate(groups)
			a.mux.Unlock()
			a.triggerSend()
		}
	}
}

func (a *accumulator) trySend(in chan<- []model.TargetGroup) {
	a.mux.Lock()
	defer a.mux.Unlock()

	select {
	case in <- a.groupsList():
		a.groupsReset()
	default:
		a.triggerSend()
	}
}

func (a *accumulator) triggerSend() {
	select {
	case a.send <- struct{}{}:
	default:
	}
}

func (a *accumulator) groupsUpdate(groups []model.TargetGroup) {
	for _, group := range groups {
		a.groups[group.Source()] = group
	}
}

func (a *accumulator) groupsReset() {
	for key := range a.groups {
		delete(a.groups, key)
	}
}

func (a *accumulator) groupsList() []model.TargetGroup {
	groups := make([]model.TargetGroup, 0, len(a.groups))
	for _, group := range a.groups {
		if group != nil {
			groups = append(groups, group)
		}
	}
	return groups
}
