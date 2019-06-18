package pihole

import (
	"github.com/netdata/go-orchestrator/module"
)

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Chart is an alias for module.Chart
	Chart = module.Chart
	// Dims is an alias for module.Dims
	Dims = module.Dims
	// Dim is an alias for module.Dim
	Dim = module.Dim
)

var charts = Charts{
	{
		ID:    "random",
		Title: "A Random Number", Units: "random", Fam: "random",
		Dims: Dims{
			{ID: "random0", Name: "random 0"},
			{ID: "random1", Name: "random 1"},
		},
	},
}

var authCharts = Charts{
	{
		ID:    "processed_dns_queries_types",
		Title: "Processed DNS Queries By Types",
		Units: "percentage",
		Fam:   "queries",
		Ctx:   "pihole.processed_dns_queries_types",
		Dims: Dims{
			{ID: "A", Div: 100},
			{ID: "AAAA", Div: 100},
			{ID: "ANY", Div: 100},
			{ID: "PTR", Div: 100},
			{ID: "SOA", Div: 100},
			{ID: "SRV", Div: 100},
			{ID: "TXT", Div: 100},
		},
	},
	{
		ID:    "forwarded_dns_queries_targets",
		Title: "Forwarded DNS Queries By Target",
		Units: "percentage",
		Fam:   "queries",
		Ctx:   "pihole.forwarded_dns_queries_targets",
		Dims: Dims{
			{ID: "target_blocklist", Name: "blocklist", Div: 100},
			{ID: "target_cache", Name: "cache", Div: 100},
		},
	},
	{
		ID:    "top_clients",
		Title: "Top Clients",
		Units: "queries",
		Fam:   "top clients",
		Ctx:   "pihole.top_clients_queries",
	},
	{
		ID:    "top_domains",
		Title: "Top Domains",
		Units: "queries",
		Fam:   "top domains",
		Ctx:   "pihole.top_domains_queries",
	},
	{
		ID:    "top_advertisers",
		Title: "Top Advertisers",
		Units: "queries",
		Fam:   "top ads",
		Ctx:   "pihole.top_ads_queries",
	},
}

func (p *Pihole) updateCharts(pmx *piholeMetrics) {
	p.updateForwardDestinationsCharts(pmx)
	p.updateTopClientChart(pmx)
	p.updateTopDomainsChart(pmx)
	p.updateTopAdvertisersChart(pmx)
}

func (p *Pihole) updateForwardDestinationsCharts(pmx *piholeMetrics) {
	if pmx.forwardDestinations == nil {
		return
	}
	chart := p.charts.Get("forwarded_dns_queries_targets")

	for _, v := range *pmx.forwardDestinations {
		if v.Name == "blocklist" || v.Name == "cache" {
			continue
		}

		id := "target_" + v.Name
		if chart.HasDim(id) {
			continue
		}

		panicIf(chart.AddDim(&Dim{ID: id, Name: v.Name, Div: 100}))
		chart.MarkNotCreated()
	}
}

func (p *Pihole) updateTopClientChart(pmx *piholeMetrics) {
	if pmx.topClients == nil {
		return
	}

	chart := p.charts.Get("top_clients")
	set := make(map[string]bool)

	for _, v := range *pmx.topClients {
		id := "top_client_" + v.Name
		set[id] = true

		if p.collected.topClients[id] {
			continue
		}

		p.collected.topClients[id] = true
		panicIf(chart.AddDim(&Dim{ID: id, Name: v.Name}))
		chart.MarkNotCreated()
	}

	removeNotUpdatedDims(chart, p.collected.topClients, set)
}

func (p *Pihole) updateTopDomainsChart(pmx *piholeMetrics) {
	if pmx.topItems == nil || pmx.topItems.TopQueries == nil {
		return
	}

	chart := p.charts.Get("top_domains")
	set := make(map[string]bool)

	for _, v := range pmx.topItems.TopQueries {
		id := "top_domain_" + v.Name
		set[id] = true

		if p.collected.topDomains[id] {
			continue
		}

		p.collected.topDomains[id] = true
		panicIf(chart.AddDim(&Dim{ID: id, Name: v.Name}))
		chart.MarkNotCreated()
	}

	removeNotUpdatedDims(chart, p.collected.topDomains, set)
}

func (p *Pihole) updateTopAdvertisersChart(pmx *piholeMetrics) {
	if pmx.topItems == nil || pmx.topItems.TopAds == nil {
		return
	}

	chart := p.charts.Get("top_advertisers")
	set := make(map[string]bool)

	for _, v := range pmx.topItems.TopAds {
		id := "top_ad_" + v.Name
		set[id] = true

		if p.collected.topAds[id] {
			continue
		}

		p.collected.topAds[id] = true
		panicIf(chart.AddDim(&Dim{ID: id, Name: v.Name}))
		chart.MarkNotCreated()
	}

	removeNotUpdatedDims(chart, p.collected.topAds, set)
}

func removeNotUpdatedDims(chart *Chart, existed, updated map[string]bool) {
	for id := range existed {
		if updated[id] {
			continue
		}

		delete(updated, id)
		panicIf(chart.RemoveDim(id))
		chart.MarkNotCreated()
	}
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
