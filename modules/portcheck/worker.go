package portcheck

import (
	"fmt"
	"net"
	"time"
)

type worker struct {
	alive    bool
	stopHook chan struct{}

	doCh        chan *port
	doneCh      chan struct{}
	host        string
	dialTimeout time.Duration
}

func newWorker(host string, dialTimeout time.Duration, doCh chan *port, doneCh chan struct{}) *worker {
	w := &worker{
		stopHook:    make(chan struct{}),
		doCh:        doCh,
		doneCh:      doneCh,
		host:        host,
		dialTimeout: dialTimeout,
		alive:       true,
	}

	go func() {
	LOOP:
		for {
			select {
			case <-w.stopHook:
				w.alive = false
				w.stopHook <- struct{}{}
				break LOOP
			case p := <-w.doCh:
				w.doWork(p)
			}
		}
	}()

	return w
}

func (w *worker) stop() {
	w.stopHook <- struct{}{}
	<-w.stopHook
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

	w.doneCh <- struct{}{}
}
