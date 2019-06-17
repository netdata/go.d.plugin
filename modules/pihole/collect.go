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
	rmx := p.getAllMetrics()
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

func (p *Pihole) getAllMetrics() *rawMetrics {
	rmx := new(rawMetrics)
	wg := &sync.WaitGroup{}

	doSummary := func() error {
		var err error
		rmx.summary, err = p.client.SummaryRaw()
		return err
	}
	doQueryTypes := func() error {
		var err error
		rmx.queryTypes, err = p.client.QueryTypes()
		return err
	}
	doForwardDestinations := func() error {
		var err error
		rmx.forwardDestinations, err = p.client.ForwardDestinations()
		return err
	}
	doTopClients := func() error {
		var err error
		rmx.topClients, err = p.client.TopClients(5)
		return err
	}
	doTopItems := func() error {
		var err error
		rmx.topItems, err = p.client.TopItems(5)
		return err
	}

	wrapper := func(do func() error) func() {
		w := func() {
			if err := do(); err != nil {
				p.Error(err)
			}
			wg.Done()
		}
		return w
	}

	tasks := []func() error{doSummary, doQueryTypes, doForwardDestinations, doTopClients, doTopItems}

	wg.Add(len(tasks))
	for _, task := range tasks {
		go wrapper(task)
	}

	wg.Wait()

	return rmx
}
