package pihole

import (
	"github.com/netdata/go.d.plugin/modules/pihole/client"

	"sync"
)

type rawMetrics struct {
	summary             *client.SummaryRaw
	queryTypes          *client.QueryTypes
	forwardDestinations *client.ForwardDestinations
	topClients          *client.TopClients
	topItems            *client.TopItems
}

func (p *Pihole) collect() (map[string]int64, error) {
	rmx := p.collectRawMetrics(true)
	mx := make(map[string]int64)

	p.collectSummary(mx, rmx)
	p.collectQueryTypes(mx, rmx)
	p.collectForwardDestination(mx, rmx)
	p.collectTopClients(mx, rmx)
	p.collectTopItems(mx, rmx)

	return mx, nil
}

func (p *Pihole) collectSummary(mx map[string]int64, rmx *rawMetrics) {
	if rmx.summary == nil {
		return
	}
}

func (p *Pihole) collectQueryTypes(mx map[string]int64, rmx *rawMetrics) {
	if rmx.queryTypes == nil {
		return
	}
}

func (p *Pihole) collectForwardDestination(mx map[string]int64, rmx *rawMetrics) {
	if rmx.forwardDestinations == nil {
		return
	}
}

func (p *Pihole) collectTopClients(mx map[string]int64, rmx *rawMetrics) {
	if rmx.topClients == nil {
		return
	}
}

func (p *Pihole) collectTopItems(mx map[string]int64, rmx *rawMetrics) {
	if rmx.topItems == nil {
		return
	}
}

func (p *Pihole) collectRawMetrics(doConcurrently bool) *rawMetrics {
	rmx := new(rawMetrics)
	wg := &sync.WaitGroup{}

	type task func() error

	logWrap := func(task task) func() {
		return func() {
			err := task()
			if err != nil {
				p.Error(err)
			}
		}
	}

	wgWrap := func(call func()) func() {
		return func() {
			call()
			wg.Done()
		}
	}

	taskSummary := func() error {
		var err error
		rmx.summary, err = p.client.SummaryRaw()
		return err
	}
	taskQueryTypes := func() error {
		var err error
		rmx.queryTypes, err = p.client.QueryTypes()
		return err
	}
	taskForwardDestinations := func() error {
		var err error
		rmx.forwardDestinations, err = p.client.ForwardDestinations()
		return err
	}
	taskTopClients := func() error {
		var err error
		rmx.topClients, err = p.client.TopClients(defaultTopClients)
		return err
	}
	taskTopItems := func() error {
		var err error
		rmx.topItems, err = p.client.TopItems(defaultTopItems)
		return err
	}

	tasks := []task{taskSummary, taskQueryTypes, taskForwardDestinations, taskTopClients, taskTopItems}

	for _, t := range tasks {
		wrapped := logWrap(t)

		if doConcurrently {
			wg.Add(1)
			go wgWrap(wrapped)
		} else {
			wrapped()
		}
	}

	wg.Wait()

	return rmx
}
