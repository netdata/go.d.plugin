// SPDX-License-Identifier: GPL-3.0-or-later

package pipeline

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
		tggs:      make(map[string]model.TargetGroup),
	}
}

type (
	accumulator struct {
		tasks     []accumulateTask
		send      chan struct{}
		sendEvery time.Duration
		mux       *sync.Mutex
		tggs      map[string]model.TargetGroup
	}
	accumulateTask struct {
		disc model.Discoverer
		tags model.Tags
	}
)

func (a *accumulator) run(ctx context.Context, in chan []model.TargetGroup) {
	var wg sync.WaitGroup
	updates := make(chan []model.TargetGroup)

	for _, task := range a.tasks {
		task := task
		wg.Add(1)
		go func() { defer wg.Done(); a.runTask(ctx, task, updates) }()

	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		tk := time.NewTicker(a.sendEvery)
		defer tk.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-tk.C:
				select {
				case <-a.send:
					a.trySend(in)
				default:
				}
			}
		}
	}()

	wg.Wait()
	<-ctx.Done()
}

func (a *accumulator) runTask(ctx context.Context, task accumulateTask, updates chan []model.TargetGroup) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() { defer wg.Done(); task.disc.Discover(ctx, updates) }()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case tggs, ok := <-updates:
			if !ok {
				break loop
			}

			for _, tgg := range tggs {
				for _, tgt := range tgg.Targets() {
					tgt.Tags().Merge(task.tags)
				}
			}

			a.mux.Lock()
			a.groupsUpdate(tggs)
			a.mux.Unlock()
			a.triggerSend()
		}
	}

	wg.Wait()
	<-ctx.Done()
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

func (a *accumulator) groupsUpdate(tggs []model.TargetGroup) {
	for _, tgg := range tggs {
		a.tggs[tgg.Source()] = tgg
	}
}

func (a *accumulator) groupsReset() {
	for key := range a.tggs {
		delete(a.tggs, key)
	}
}

func (a *accumulator) groupsList() []model.TargetGroup {
	tggs := make([]model.TargetGroup, 0, len(a.tggs))
	for _, tgg := range a.tggs {
		if tgg != nil {
			tggs = append(tggs, tgg)
		}
	}
	return tggs
}
