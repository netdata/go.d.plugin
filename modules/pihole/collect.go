// SPDX-License-Identifier: GPL-3.0-or-later

package pihole

import (
	"sync"
	"time"
)

func (p *Pihole) collect() (map[string]int64, error) {
	pmx := new(piholeMetrics)
	p.scrapePihole(pmx, true)
	mx := make(map[string]int64)

	// non auth
	collectSummary(mx, pmx)
	// auth
	collectQueryTypes(mx, pmx)
	collectForwardDestination(mx, pmx)
	collectTopClients(mx, pmx)
	collectTopDomains(mx, pmx)

	p.updateCharts(pmx)

	return mx, nil
}

func collectSummary(mx map[string]int64, pmx *piholeMetrics) {
	if !pmx.hasSummary() {
		return
	}

	mx["ads_blocked_today"] = pmx.summary.AdsBlockedToday
	mx["ads_percentage_today"] = int64(pmx.summary.AdsPercentageToday * 100)
	mx["domains_being_blocked"] = pmx.summary.DomainsBeingBlocked
	// GravityLastUpdated.Absolute is <nil> if the file is not exists (deleted/moved)
	if pmx.summary.GravityLastUpdated.Absolute != nil {
		mx["blocklist_last_update"] = time.Now().Unix() - *pmx.summary.GravityLastUpdated.Absolute
	}
	mx["dns_queries_today"] = pmx.summary.DNSQueriesToday
	mx["queries_forwarded"] = pmx.summary.QueriesForwarded
	mx["queries_cached"] = pmx.summary.QueriesCached
	mx["unique_clients"] = pmx.summary.UniqueClients
	mx["file_exists"] = boolToInt(pmx.summary.GravityLastUpdated.FileExists)
	mx["status"] = boolToInt(pmx.summary.Status == "enabled")
}

func collectQueryTypes(mx map[string]int64, pmx *piholeMetrics) {
	if !pmx.hasQueryTypes() {
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
	if !pmx.hasForwardDestinations() {
		return
	}
	for _, v := range *pmx.forwarders {
		mx["destination_"+v.Name] = int64(v.Percent * 100)
	}
}

func collectTopClients(mx map[string]int64, pmx *piholeMetrics) {
	if !pmx.hasTopClients() {
		return
	}
	for _, v := range *pmx.topClients {
		mx["top_client_"+v.Name] = v.Requests
	}
}

func collectTopDomains(mx map[string]int64, pmx *piholeMetrics) {
	if pmx.hasTopQueries() {
		for _, v := range *pmx.topQueries {
			mx["top_perm_domain_"+v.Name] = v.Hits
		}
	}
	if pmx.hasTopAdvertisers() {
		for _, v := range *pmx.topAds {
			mx["top_blocked_domain_"+v.Name] = v.Hits
		}
	}
}

func (p *Pihole) scrapePihole(pmx *piholeMetrics, doConcurrently bool) {
	type task func(*piholeMetrics)

	var tasks = []task{p.scrapeSummary}

	if p.Password != "" {
		tasks = []task{
			p.scrapeSummary,
			p.scrapeQueryTypes,
			p.scrapeForwardedDestinations,
			p.scrapeTopClients,
			p.scrapeTopItems,
		}
	}

	wg := &sync.WaitGroup{}

	wrap := func(call task) task {
		return func(metrics *piholeMetrics) {
			call(metrics)
			wg.Done()
		}
	}

	for _, task := range tasks {
		if doConcurrently {
			wg.Add(1)
			task = wrap(task)
			go task(pmx)
		} else {
			task(pmx)
		}
	}

	wg.Wait()
}

func (p *Pihole) scrapeSummary(pmx *piholeMetrics) {
	v, err := p.client.SummaryRaw()
	if err != nil {
		p.Error(err)
		return
	}
	pmx.summary = v
}

func (p *Pihole) scrapeQueryTypes(pmx *piholeMetrics) {
	v, err := p.client.QueryTypes()
	if err != nil {
		p.Error(err)
		return
	}
	pmx.queryTypes = v
}

func (p *Pihole) scrapeForwardedDestinations(pmx *piholeMetrics) {
	v, err := p.client.ForwardDestinations()
	if err != nil {
		p.Error(err)
		return
	}
	pmx.forwarders = v
}

func (p *Pihole) scrapeTopClients(pmx *piholeMetrics) {
	v, err := p.client.TopClients(p.TopClientsEntries)
	if err != nil {
		p.Error(err)
		return
	}
	pmx.topClients = v
}

func (p *Pihole) scrapeTopItems(pmx *piholeMetrics) {
	v, err := p.client.TopItems(p.TopItemsEntries)
	if err != nil {
		p.Error(err)
		return
	}
	pmx.topQueries = &v.TopQueries
	pmx.topAds = &v.TopAds
}

func boolToInt(b bool) int64 {
	if !b {
		return 0
	}
	return 1
}
