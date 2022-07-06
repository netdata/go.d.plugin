// SPDX-License-Identifier: GPL-3.0-or-later

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
	w := &worker{
		task:      task,
		taskDone:  taskDone,
		exchanger: exchanger,
	}

	go w.workLoop()

	return w
}

type worker struct {
	task     chan task
	taskDone chan struct{}

	exchanger exchanger
}

func (w *worker) workLoop() {
	for task := range w.task {
		w.doWork(task)
	}
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
