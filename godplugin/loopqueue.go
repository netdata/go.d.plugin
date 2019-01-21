package godplugin

import "sync"

type loopQueue struct {
	mux   sync.Mutex
	queue []Job
}

func (q *loopQueue) len() int {
	q.mux.Lock()
	defer q.mux.Unlock()
	return len(q.queue)
}

func (q *loopQueue) add(job Job) {
	q.mux.Lock()
	defer q.mux.Unlock()

	q.queue = append(q.queue, job)
}

func (q *loopQueue) remove(fullName string) Job {
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

func (q *loopQueue) notify(clock int) {
	q.mux.Lock()
	defer q.mux.Unlock()

	for _, job := range q.queue {
		job.Tick(clock)
	}
}
