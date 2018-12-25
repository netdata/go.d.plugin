package portcheck

import (
	"fmt"
	"net"
	"time"
)

type worker struct {
	task     chan *port
	taskDone chan struct{}

	host        string
	dialTimeout time.Duration
}

func newWorker(host string, dialTimeout time.Duration, task chan *port, taskDone chan struct{}) *worker {
	w := &worker{
		task:        task,
		taskDone:    taskDone,
		host:        host,
		dialTimeout: dialTimeout,
	}

	go w.workLoop()

	return w
}

func (w *worker) workLoop() {
	for task := range w.task {
		w.doWork(task)
	}
}

func (w *worker) doWork(port *port) {
	t := time.Now()
	c, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", w.host, port.number), w.dialTimeout)
	port.latency = time.Since(t)

	if err == nil {
		port.setState(success)
		_ = c.Close()
	}

	if err != nil {
		v, ok := err.(interface{ Timeout() bool })

		if ok && v.Timeout() {
			port.setState(timeout)
		} else {
			port.setState(failed)
		}
	}

	w.taskDone <- struct{}{}
}
