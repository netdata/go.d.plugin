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
	rmx := &rawMetrics{}
	var wg sync.WaitGroup

	wg.Add(5)

	go func() {
		var err error
		rmx.summary, err = p.client.SummaryRaw()
		if err != nil {
			p.Error(err)
		}
		wg.Done()
	}()
	go func() {
		var err error
		rmx.queryTypes, err = p.client.QueryTypes()
		if err != nil {
			p.Error(err)
		}
		wg.Done()
	}()
	go func() {
		var err error
		rmx.forwardDestinations, err = p.client.ForwardDestinations()
		if err != nil {
			p.Error(err)
		}
		wg.Done()
	}()
	go func() {
		var err error
		rmx.topClients, err = p.client.TopClients(5)
		if err != nil {
			p.Error(err)
		}
		wg.Done()
	}()
	go func() {
		var err error
		rmx.topItems, err = p.client.TopItems(5)
		if err != nil {
			p.Error(err)
		}
		wg.Done()
	}()

	wg.Wait()

	return rmx
}
