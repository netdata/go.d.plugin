package run

import (
	"context"
	"sync"
	"time"

	jobpkg "github.com/netdata/go.d.plugin/agent/job"
	"github.com/netdata/go.d.plugin/agent/ticker"
	"github.com/netdata/go.d.plugin/logger"
)

type (
	Manager struct {
		mux   sync.Mutex
		queue queue
		*logger.Logger
	}
	queue []jobpkg.Job
)

func NewManager() *Manager {
	return &Manager{
		mux:    sync.Mutex{},
		Logger: logger.New("run", "manager"),
	}
}

func (m *Manager) Run(ctx context.Context) {
	m.Info("instance is started")
	defer func() { m.Info("instance is stopped") }()

	tk := ticker.New(time.Second)
	defer tk.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case clock := <-tk.C:
			m.Debugf("tick %d", clock)
			m.notify(clock)
		}
	}
}

// Starts starts a job and adds it to the job queue.
func (m *Manager) Start(job jobpkg.Job) {
	m.mux.Lock()
	defer m.mux.Unlock()

	go job.Start()
	m.queue.add(job)
}

// Stop stops a job and removes it from the job queue.
func (m *Manager) Stop(fullName string) {
	m.mux.Lock()
	defer m.mux.Unlock()

	if job := m.queue.remove(fullName); job != nil {
		job.Stop()
	}
}

// Cleanup stops all jobs in the queue.
func (m *Manager) Cleanup() {
	for _, v := range m.queue {
		v.Stop()
	}
	m.queue = m.queue[:0]
}

func (m *Manager) notify(clock int) {
	m.mux.Lock()
	defer m.mux.Unlock()

	for _, v := range m.queue {
		v.Tick(clock)
	}
}

func (q *queue) add(job jobpkg.Job) {
	*q = append(*q, job)
}

func (q *queue) remove(fullName string) jobpkg.Job {
	for idx, v := range *q {
		if v.FullName() != fullName {
			continue
		}
		j := (*q)[idx]
		copy((*q)[idx:], (*q)[idx+1:])
		(*q)[len(*q)-1] = nil
		*q = (*q)[:len(*q)-1]
		return j
	}
	return nil
}
