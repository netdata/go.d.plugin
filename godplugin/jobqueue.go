package godplugin

import "sync"

type jobQueue struct {
	mux   sync.Mutex
	queue []Job
}

func (q *jobQueue) add(job Job) {
	q.mux.Lock()
	defer q.mux.Unlock()

	q.queue = append(q.queue, job)
}

func (q *jobQueue) remove(fullName string) Job {
	q.mux.Lock()
	defer q.mux.Unlock()

	for i, job := range q.queue {
		if job.FullName() == fullName {
			q.queue = append(q.queue[:i], q.queue[i+1:]...)
			return job
		}
	}
	return nil
}

func (q *jobQueue) notify(clock int) {
	q.mux.Lock()
	defer q.mux.Unlock()

	for _, job := range q.queue {
		job.Tick(clock)
	}
}
