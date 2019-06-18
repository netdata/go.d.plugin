package pihole

import (
	"github.com/netdata/go.d.plugin/modules/pihole/client"

	"sync"
)

type piholeMetrics struct {
	summary             *client.SummaryRaw
	queryTypes          *client.QueryTypes
	forwardDestinations *client.ForwardDestinations
	topClients          *client.TopClients
	topItems            *client.TopItems
}

func (p *Pihole) collect() (map[string]int64, error) {
	pmx := p.scrapePihole(true)
	mx := make(map[string]int64)

	// non auth
	collectSummary(mx, pmx)
	// auth
	collectQueryTypes(mx, pmx)
	collectForwardDestination(mx, pmx)
	collectTopClients(mx, pmx)
	collectTopItems(mx, pmx)

	p.updateCharts(pmx)

	return mx, nil
}

func collectSummary(mx map[string]int64, pmx *piholeMetrics) {
	if pmx.summary == nil {
		return
	}
}

func collectQueryTypes(mx map[string]int64, pmx *piholeMetrics) {
	if pmx.queryTypes == nil {
		return
	}
	mx["A"] = int64(pmx.queryTypes.A * 100)
	mx["AAAA"] = int64(pmx.queryTypes.AAAA * 100)
	mx["ANY"] = int64(pmx.queryTypes.ANY * 100)
	mx["PTR"] = int64(pmx.queryTypes.PTR * 100)
	mx["SOA"] = int64(pmx.queryTypes.SOA * 100)
	mx["SRV"] = int64(pmx.queryTypes.SRV * 100)
	mx["TXT"] = int64(pmx.queryTypes.TXT * 100)
}

func collectForwardDestination(mx map[string]int64, pmx *piholeMetrics) {
	if pmx.forwardDestinations == nil {
		return
	}
	for _, v := range *pmx.forwardDestinations {
		mx["target_"+v.Name] = int64(v.Percent * 100)
	}
}

func collectTopClients(mx map[string]int64, pmx *piholeMetrics) {
	if pmx.topItems == nil {
		return
	}
	for _, v := range *pmx.topClients {
		mx["top_client_"+v.Name] = v.Queries
	}
}

func collectTopItems(mx map[string]int64, pmx *piholeMetrics) {
	if pmx.topItems == nil {
		return
	}
	for _, v := range pmx.topItems.TopQueries {
		mx["top_domain_"+v.Name] = v.Queries
	}
	for _, v := range pmx.topItems.TopAds {
		mx["top_ad_"+v.Name] = v.Queries
	}
}

func (p *Pihole) scrapePihole(doConcurrently bool) *piholeMetrics {
	pmx := new(piholeMetrics)

	taskSummary := func() {
		var err error
		pmx.summary, err = p.client.SummaryRaw()
		if err != nil {
			p.Error(err)
		}
	}
	taskQueryTypes := func() {
		var err error
		pmx.queryTypes, err = p.client.QueryTypes()
		if err != nil {
			p.Error(err)
		}
	}
	taskForwardDestinations := func() {
		var err error
		pmx.forwardDestinations, err = p.client.ForwardDestinations()
		if err != nil {
			p.Error(err)
		}
	}
	taskTopClients := func() {
		var err error
		pmx.topClients, err = p.client.TopClients(defaultTopClients)
		if err != nil {
			p.Error(err)
		}
	}
	taskTopItems := func() {
		var err error
		pmx.topItems, err = p.client.TopItems(defaultTopItems)
		if err != nil {
			p.Error(err)
		}
	}

	var tasks = []func(){taskSummary}
	if p.Password != "" {
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

	return pmx
}
