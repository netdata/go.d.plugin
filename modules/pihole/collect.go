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

func (r rawMetrics) hasSummary() bool { return r.summary != nil }

func (r rawMetrics) hasQueryTypes() bool { return r.queryTypes != nil }

func (r rawMetrics) hasForwardDestinations() bool { return r.forwardDestinations != nil }

func (r rawMetrics) hasTopClients() bool { return r.topClients != nil }

func (r rawMetrics) hasTopItems() bool { return r.topItems != nil }

func (p *Pihole) collect() (map[string]int64, error) {
	rmx := p.collectRawMetrics(true)
	mx := make(map[string]int64)

	// non auth
	p.collectSummary(mx, rmx)

	// auth
	p.collectQueryTypes(mx, rmx)
	p.collectForwardDestination(mx, rmx)
	p.collectTopClients(mx, rmx)
	p.collectTopItems(mx, rmx)

	p.updateCharts(rmx)

	return mx, nil
}

func (p *Pihole) collectSummary(mx map[string]int64, rmx *rawMetrics) {
	if !rmx.hasSummary() {
		return
	}
}

func (p *Pihole) collectQueryTypes(mx map[string]int64, rmx *rawMetrics) {
	if rmx.hasQueryTypes() {
		return
	}
	mx["A"] = int64(rmx.queryTypes.A * 100)
	mx["AAAA"] = int64(rmx.queryTypes.AAAA * 100)
	mx["ANY"] = int64(rmx.queryTypes.ANY * 100)
	mx["PTR"] = int64(rmx.queryTypes.PTR * 100)
	mx["SOA"] = int64(rmx.queryTypes.SOA * 100)
	mx["SRV"] = int64(rmx.queryTypes.SRV * 100)
	mx["TXT"] = int64(rmx.queryTypes.TXT * 100)
}

func (p *Pihole) collectForwardDestination(mx map[string]int64, rmx *rawMetrics) {
	if rmx.hasForwardDestinations() {
		return
	}
	for _, v := range *rmx.forwardDestinations {
		mx["target_"+v.Name] = int64(v.Percent * 100)
	}
}

func (p *Pihole) collectTopClients(mx map[string]int64, rmx *rawMetrics) {
	if rmx.hasTopClients() {
		return
	}
	for _, v := range *rmx.topClients {
		mx["top_client_"+v.Name] = v.Queries
	}
}

func (p *Pihole) collectTopItems(mx map[string]int64, rmx *rawMetrics) {
	if rmx.hasTopItems() {
		return
	}
	for _, v := range rmx.topItems.TopQueries {
		mx["top_query_"+v.Name] = v.Queries
	}
	for _, v := range rmx.topItems.TopAds {
		mx["top_ad_"+v.Name] = v.Queries
	}
}

func (p *Pihole) collectRawMetrics(doConcurrently bool) *rawMetrics {
	rmx := new(rawMetrics)

	taskSummary := func() {
		var err error
		rmx.summary, err = p.client.SummaryRaw()
		if err != nil {
			p.Error(err)
		}
	}
	taskQueryTypes := func() {
		var err error
		rmx.queryTypes, err = p.client.QueryTypes()
		if err != nil {
			p.Error(err)
		}
	}
	taskForwardDestinations := func() {
		var err error
		rmx.forwardDestinations, err = p.client.ForwardDestinations()
		if err != nil {
			p.Error(err)
		}
	}
	taskTopClients := func() {
		var err error
		rmx.topClients, err = p.client.TopClients(defaultTopClients)
		if err != nil {
			p.Error(err)
		}
	}
	taskTopItems := func() {
		var err error
		rmx.topItems, err = p.client.TopItems(defaultTopItems)
		if err != nil {
			p.Error(err)
		}
	}

	var tasks = []func(){taskSummary}
	if p.client.WebPassword != "" {
		tasks = []func(){
			taskSummary,
			taskQueryTypes,
			taskForwardDestinations,
			taskTopClients,
			taskTopItems,
		}
	}

	wg := &sync.WaitGroup{}

	wrap := func(call func()) func() {
		return func() {
			call()
			wg.Done()
		}
	}

	for _, task := range tasks {
		if doConcurrently {
			wg.Add(1)
			task = wrap(task)
			go task()
		} else {
			task()
		}
	}

	wg.Wait()

	return rmx
}
