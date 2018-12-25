package dnsquery

import (
	"net"
	"strconv"

	"github.com/miekg/dns"
)

type task struct {
	server *server
	domain string
	rtype  uint16
}

func newWorker(exchanger exchanger, task chan task, taskDone chan struct{}) *worker {
	return &worker{
		shutdown:  make(chan struct{}),
		task:      task,
		taskDone:  taskDone,
		exchanger: exchanger,
	}
}

type worker struct {
	shutdown chan struct{}
	task     chan task
	taskDone chan struct{}

	exchanger exchanger
}

func (w *worker) workLoop() {
LOOP:
	for {
		select {
		case <-w.shutdown:
			w.shutdown <- struct{}{}
			break LOOP
		case task := <-w.task:
			w.doWork(task)
		}
	}
}

func (w *worker) stop() {
	w.shutdown <- struct{}{}
	<-w.shutdown
}

func (w *worker) doWork(t task) {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(t.domain), t.rtype)
	address := net.JoinHostPort(t.server.name, strconv.Itoa(t.server.port))

	resp, rtt, err := w.exchanger.Exchange(msg, address)

	t.server.resp = resp
	t.server.rtt = rtt
	t.server.err = err

	w.taskDone <- struct{}{}
}
