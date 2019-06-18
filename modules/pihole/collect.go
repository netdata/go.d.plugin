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

	taskSummary := func() {
		var err error
		if rmx.summary, err = p.client.SummaryRaw(); err != nil {
			p.Error(err)
		}
	}
	taskQueryTypes := func() {
		var err error
		if rmx.queryTypes, err = p.client.QueryTypes(); err != nil {
			p.Error(err)
		}
	}
	taskForwardDestinations := func() {
		var err error
		if rmx.forwardDestinations, err = p.client.ForwardDestinations(); err != nil {
			p.Error(err)
		}
	}
	taskTopClients := func() {
		var err error
		if rmx.topClients, err = p.client.TopClients(defaultTopClients); err != nil {
			p.Error(err)
		}
	}
	taskTopItems := func() {
		var err error
		if rmx.topItems, err = p.client.TopItems(defaultTopItems); err != nil {
			p.Error(err)
		}
	}

	wg := &sync.WaitGroup{}

	wrapper := func(call func()) func() {
		return func() {
			call()
			wg.Done()
		}
	}

	tasks := []func(){taskSummary, taskQueryTypes, taskForwardDestinations, taskTopClients, taskTopItems}

	for _, task := range tasks {
		if doConcurrently {
			wg.Add(1)
			go wrapper(task)()
		} else {
			task()
		}
	}

	wg.Wait()

	return rmx
}
