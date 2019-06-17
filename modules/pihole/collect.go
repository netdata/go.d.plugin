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

	p.collectSummary(mx, rmx.summary)
	p.collectQueryTypes(mx, rmx.queryTypes)
	p.collectForwardDestination(mx, rmx.forwardDestinations)
	p.collectTopClients(mx, rmx.topClients)
	p.collectTopItems(mx, rmx.topItems)

	return mx, nil
}

func (p *Pihole) collectSummary(mx map[string]int64, summary *client.SummaryRaw) {
	if summary == nil {
		return
	}
}

func (p *Pihole) collectQueryTypes(mx map[string]int64, queryTypes *client.QueryTypes) {
	if queryTypes == nil {
		return
	}
}

func (p *Pihole) collectForwardDestination(mx map[string]int64, forwardDest *client.ForwardDestinations) {
	if forwardDest == nil {
		return
	}
}

func (p *Pihole) collectTopClients(mx map[string]int64, topClients *client.TopClients) {
	if topClients == nil {
		return
	}
}

func (p *Pihole) collectTopItems(mx map[string]int64, topItems *client.TopItems) {
	if topItems == nil || topItems.TopAds == nil && topItems.TopQueries == nil {
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
